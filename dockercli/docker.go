package dockercli

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"github.com/moonlightming/simple-docker-inside-webhook/commons"
	"github.com/moonlightming/simple-docker-inside-webhook/conf"
	"io"
	"log"
	"os"
	"time"
)

var (
	config     = conf.NewConfig() // Global config
	cli        = newDockerCli()   // single Docker client
	authBase64 = ""
)

func init() {
	// Auth for base64
	if config.DockerRegistryAuth.User != "" && config.DockerRegistryAuth.Password != "" {
		auth := types.AuthConfig{
			Username: config.DockerRegistryAuth.User,
			Password: config.DockerRegistryAuth.Password,
		}
		authBytes, err := json.Marshal(auth)
		if err != nil {
			panic(err)
		}
		authBase64 = base64.URLEncoding.EncodeToString(authBytes)
	}
}

func newDockerCli() *client.Client {
	cli, err := client.NewClientWithOpts(
		client.WithHost(fmt.Sprintf("tcp://%s", config.DockerHost)),
		client.WithVersion(config.DockerApiVersion),
	)
	if err != nil {
		panic(err)
	}
	return cli
}

func PullImage(imageFullName string, repoType string) error {
	reader, err := cli.ImagePull(
		context.Background(),
		imageFullName,
		types.ImagePullOptions{RegistryAuth: isPrivate(repoType)},
	)
	if err != nil {
		return err
	}
	io.Copy(os.Stdout, reader)
	return nil
}

func ListServiceWithName(serviceName string) ([]swarm.Service, error) {
	var (
		swarms []swarm.Service
		err    error
	)
	if swarms, err = cli.ServiceList(
		context.Background(),
		types.ServiceListOptions{
			Filters: filters.NewArgs(filters.Arg("name", serviceName)),
		},
	); err != nil {
		return nil, err
	}
	return swarms, nil
}

func UpdateService(hookRequest commons.HookRequest, groupName string) error {
	var beUpService swarm.Service
	if swarms, err := ListServiceWithName(groupName + "_" + hookRequest.Name); err != nil {
		return err
	} else {
		beUpService = swarms[0]
	}
	beUpServiceJ, err := json.Marshal(beUpService)
	log.Printf("BeforeInspectService: %s", beUpServiceJ)
	log.Println("############## Update #####################")

	// if you want update Service for new Image , must change the Image field like this
	// the Image field likes "{ImagesFullName}@{Images_PushData_Digest}"
	// Example: beUpService.Spec.TaskTemplate.ContainerSpec.Image = "registry.cn-shenzhen.aliyuncs.com/moonlightming/xblog-hugo:latest@sha256:75e8d7e28743402ec93dfa05cdc45c12c920e59fd98084eb3cf65615f955c5f9"
	beUpService.Spec.TaskTemplate.ContainerSpec.Image = fmt.Sprintf("%s@%s", hookRequest.PublicAddr(), hookRequest.PushData.Digest)

	// Update Service
	warning, err := cli.ServiceUpdate(
		context.Background(),
		beUpService.ID,
		swarm.Version{Index: beUpService.Version.Index},
		beUpService.Spec,
		types.ServiceUpdateOptions{EncodedRegistryAuth: isPrivate(hookRequest.RepoType), QueryRegistry: false},
	)
	log.Printf("Warning: %+v", warning)
	log.Printf("Err: %+v", err)
	return err
}

func CleanNoneTagImage() error {
	var (
		images []types.ImageSummary
		err    error
	)
	if images, err = cli.ImageList(context.Background(), types.ImageListOptions{}); err != nil {
		return err
	}
	go func() {
		for _, image := range images {
			if len(image.RepoTags) == 0 || image.RepoTags[0] == "<none>:<none>" {
				log.Printf("The image will be remove: %s\n", image.ID)
				if err := RemoveImage(image.ID, true); err != nil {
					log.Printf("The image remove error: %s", err)
				}
				time.Sleep(200 * time.Millisecond)
			}
		}
	}()
	return nil
}

func RemoveImage(imageId string, force bool) error {
	if _, err := cli.ImageRemove(context.Background(), imageId, types.ImageRemoveOptions{Force: force}); err != nil {
		return err
	}
	return nil
}

// if the registry private, return auth code
func isPrivate(repoType string) string {
	if repoType == "PRIVATE" {
		return authBase64
	}
	return ""
}
