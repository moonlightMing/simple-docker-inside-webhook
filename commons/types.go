package commons

import "fmt"

const DockerServerHost = "aliyuncs.com"

type HookRequest struct {
	PushData   `json:"push_data"`
	Repository `json:"repository"`
}

type PushData struct {
	Digest   string `json:"digest"`
	PushedAt string `json:"pushed_at"`
	Tag      string `json:"tag"`
}

type Repository struct {
	DateCreated            string `json:"date_created"`
	Name                   string `json:"name"`
	Namespace              string `json:"namespace"`
	Region                 string `json:"region"`
	RepoAuthenticationType string `json:"repo_authentication_type"`
	RepoFullName           string `json:"repo_full_name"`
	RepoOriginType         string `json:"repo_origin_type"`
	RepoType               string `json:"repo_type"`
}

func (h *HookRequest) PublicAddr() (repositoryAddr string) {
	// like this:
	// 	 registry.cn-shenzhen.aliyuncs.com/moonlightming/xblog-hugo
	return fmt.Sprintf("registry.%s.%s/%s", h.Region, DockerServerHost, h.RepoFullName)
}

func (h *HookRequest) VpcAddr() (repositoryAddr string) {
	// like this:
	//   registry-vpc.cn-shenzhen.aliyuncs.com/moonlightming/xblog-hugo
	return fmt.Sprintf("registry-vpc.%s.%s/%s", h.Region, DockerServerHost, h.RepoFullName)
}

func (h *HookRequest) InternalAddr() (repositoryAddr string) {
	// like this:
	//   registry-internal.cn-shenzhen.aliyuncs.com/moonlightming/xblog-hugo
	return fmt.Sprintf("registry-internal.%s.%s/%s", h.Region, DockerServerHost, h.RepoFullName)
}
