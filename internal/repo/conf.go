package repo

import (
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

type AuthType string

const AuthTypeGithub AuthType = "github"

type Conf struct {
	Repositories []RepositoryConf `json:"repositories"`
	AuthConf     *AuthConf        `json:"auth_conf"`
}

type AuthConf struct {
	Username       string          `json:"username"`
	Email          string          `json:"email"`
	GithubAuthConf *GithubAuthConf `json:"github_auth_conf"`
}

func (a AuthConf) Auth(authType AuthType) transport.AuthMethod {
	if authType == AuthTypeGithub && a.GithubAuthConf != nil {
		return a.GithubAuthConf.Auth()
	}

	return nil
}

type GithubAuthConf struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// IsEmpty Github 鉴权信息是否为空
func (ac GithubAuthConf) IsEmpty() bool {
	return ac.Username == "" && ac.Password == ""
}

// Auth 转换为授权请求信息
func (ac GithubAuthConf) Auth() transport.AuthMethod {
	if ac.IsEmpty() {
		return nil
	}

	return &http.BasicAuth{
		Username: ac.Username,
		Password: ac.Password,
	}
}

type RepositoryConf struct {
	Name       string    `json:"name"`
	Branch     string    `json:"branch"`
	URL        string    `json:"url"`
	StorageDir string    `json:"storage_dir"`
	Type       AuthType  `json:"type"`
	AuthConf   *AuthConf `json:"auth_conf"`
	Readonly   bool      `json:"readonly"`
}

func (rc *RepositoryConf) Init(repoConf *Conf) {
	if rc.AuthConf == nil {
		rc.AuthConf = &AuthConf{}
	}

	rc.AuthConf.Merge(repoConf.AuthConf)
}

// Merge 与其他的 AuthConf 合并
func (a *AuthConf) Merge(other *AuthConf) *AuthConf {
	if other == nil || (a.GithubAuthConf != nil && !a.GithubAuthConf.IsEmpty()) {
		return a
	}

	a.GithubAuthConf = other.GithubAuthConf
	return a
}
