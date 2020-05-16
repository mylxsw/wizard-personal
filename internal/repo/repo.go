package repo

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/container"
	"github.com/mylxsw/wizard-personal/configs"
	"github.com/mylxsw/wizard-personal/internal/utils"
	"github.com/mylxsw/wizard-personal/pkg/array"
	"github.com/mylxsw/wizard-personal/pkg/dtree"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// ErrRepositoryNotExist 仓库不存在错误
var ErrRepositoryNotExist = errors.New("repository is not exist")

// ErrRepositoryNameExisted 同名的仓库已经存在
var ErrRepositoryNameExisted = errors.New("repository with same name already existed")

// Repo Git 仓库
type Repo struct {
	cc   container.Container
	conf *configs.Config
}

// NewRepo 创建 Repo 实例
func NewRepo(cc container.Container) *Repo {
	return &Repo{cc: cc, conf: configs.Get(cc)}
}

// RemoveRepository 删除仓库配置
func (repo Repo) RemoveRepository(name string) error {
	conf, err := repo.LoadConf()
	if err != nil {
		return errors.Wrap(err, "load config failed")
	}

	deleteIndex := -1
	for i, rep := range conf.Repositories {
		if rep.Name == name {
			deleteIndex = i
			break
		}
	}

	if deleteIndex >= 0 {
		conf.Repositories = append(conf.Repositories[:deleteIndex], conf.Repositories[deleteIndex+1:]...)
	}

	if err := utils.WriteJSONFile(repo.configFilePath(), conf); err != nil {
		return errors.Wrap(err, "写入配置文件失败")
	}

	return nil
}

// AddRepository 新增仓库配置
func (repo Repo) AddRepository(r RepositoryConf) error {
	conf, err := repo.LoadConf()
	if err != nil {
		return errors.Wrap(err, "load config failed")
	}

	for _, rep := range conf.Repositories {
		if rep.Name == r.Name {
			return ErrRepositoryNameExisted
		}
	}

	r.StorageDir, err = repo.buildRepositoryStorageDir(r.URL)
	if err == nil {
		return err
	}
	if r.Branch == "" {
		r.Branch = "master"
	}

	conf.Repositories = append(conf.Repositories, r)
	if err := utils.WriteJSONFile(repo.configFilePath(), conf); err != nil {
		return errors.Wrap(err, "写入配置文件失败")
	}

	return nil
}

// configFilePah 返回配置文件路径
func (repo Repo) configFilePath() string {
	return filepath.Join(repo.conf.WorkDir, ".wizard-personal.json")
}

// LoadConf 加载仓库配置
func (repo Repo) LoadConf() (*Conf, error) {
	configFilePath := repo.configFilePath()
	exist, err := utils.FileExist(configFilePath)
	if err != nil {
		return nil, err
	}

	if !exist {
		if err := utils.WriteJSONFile(configFilePath, Conf{
			Repositories: []RepositoryConf{
				{
					Name:       "default",
					Branch:     "master",
					URL:        "https://github.com/mylxsw/growing-up.git",
					StorageDir: "growing-up",
					Type:       "github",
					Readonly:   true,
				},
			},
			AuthConf: &AuthConf{
				Username: "mylxsw",
				Email:    "mylxsw@aicode.cc",
				GithubAuthConf: &GithubAuthConf{
					Username: "",
					Password: "",
				},
			},
		}); err != nil {
			return nil, err
		}
	}

	var repoConf Conf
	if err := utils.LoadJSONFile(configFilePath, &repoConf); err != nil {
		return nil, err
	}

	return &repoConf, nil
}

// GetRepoConfByName 通过仓库名获取仓库配置
func (repo Repo) GetRepoConfByName(name string) (*RepositoryConf, error) {
	return repo.GetRepoConfByCB(func(r RepositoryConf) bool {
		return r.Name == name
	})
}

// GetRepoConfByCB 通过 cb 回调函数来获取仓库配置
func (repo Repo) GetRepoConfByCB(cb func(r RepositoryConf) bool) (*RepositoryConf, error) {
	repoConf, err := repo.LoadConf()
	if err != nil {
		return nil, err
	}

	for _, r := range repoConf.Repositories {
		if cb(r) {
			r.Init(repoConf)
			return &r, nil
		}
	}

	return nil, ErrRepositoryNotExist
}

// Inspection 输出仓库信息
func (repo Repo) Inspection(gitRepo *git.Repository) {
	if gitRepo == nil {
		return
	}

	ref, err := gitRepo.Head()
	if err != nil {
		log.Errorf("获取仓库 HEAD 失败： %v", err)
		return
	}

	cIter, _ := gitRepo.Log(&git.LogOptions{From: ref.Hash()})
	if err := cIter.ForEach(func(c *object.Commit) error {
		log.Debugf("%s > %s", c.Author, c.Message)
		return nil
	}); err != nil {
		log.Errorf("遍历提交记录失败: %v", err)
	}
}

// SaveFile 保存文档
func (repo Repo) SaveFile(name string, filename string, original string, data []byte) error {
	changed := make([]string, 0)
	_, err := repo.workInWorktree(name, func(fs afero.Fs, conf *RepositoryConf) error {
		filename = filepath.Join(conf.StorageDir, filename)
		originalFilename := filepath.Join(conf.StorageDir, original)

		if err := repo.ensureDirExist(fs, filename); err != nil {
			return err
		}

		if original != "" {
			originalFileExist, err := afero.Exists(fs, originalFilename)
			if err != nil {
				return errors.Wrap(err, "检查原始文件是否存在失败")
			}

			// 原始文件名和新的文件名不一致，执行 rename 操作
			if originalFileExist && originalFilename != filename {
				if err := fs.Rename(originalFilename, filename); err != nil {
					return errors.Wrap(err, "文件重命名失败")
				}
				changed = append(changed, originalFilename, filename)
			}
		}

		existFile, err := afero.Exists(fs, filename)
		if err != nil {
			return errors.Wrap(err, "检查文件是否存在失败")
		}

		fileChanged := true
		if existFile {
			// 文件存在，则先校验文件内容是否发生变化，没有变化就直接返回成功了
			original, err := afero.ReadFile(fs, filename)
			if err != nil {
				return errors.Wrap(err, "读取原始文件失败")
			}

			if bytes.Equal(original, data) {
				fileChanged = false
			}
		}

		if fileChanged {
			if err := afero.WriteFile(fs, filename, data, os.ModePerm); err != nil {
				return errors.Wrap(err, "保存文件失败")
			}

			changed = append(changed, filename)
		}

		return nil
	})

	if err != nil {
		return err
	}

	if len(changed) > 0 {
		if err := repo.commitChanges(name, changed); err != nil {
			return err
		}
	}

	return nil
}

// DeleteFile 删除文档
func (repo Repo) DeleteFile(name string, filename string) error {
	_, err := repo.workInWorktree(name, func(fs afero.Fs, conf *RepositoryConf) error {
		filename = filepath.Join(conf.StorageDir, filename)

		exist, err := afero.Exists(fs, filename)
		if err != nil {
			return errors.Wrap(err, "文件状态查询失败")
		}

		if !exist {
			return nil
		}

		return fs.Remove(filename)
	})

	if err == nil {
		return repo.commitChanges(name, []string{filename})
	}

	return err
}

func (repo Repo) workInWorktree(name string, cb interface{}) ([]interface{}, error) {
	confProvider := func() (*RepositoryConf, error) {
		return repo.GetRepoConfByName(name)
	}

	provider, err := repo.cc.ServiceProvider(confProvider)
	if err != nil {
		return nil, err
	}

	return repo.cc.CallWithProvider(cb, provider)
}

// ensureDirExist 确保目录存在
func (repo Repo) ensureDirExist(fs afero.Fs, filename string) error {
	existDir, err := afero.DirExists(fs, filepath.Dir(filename))
	if err != nil {
		return errors.Wrap(err, "检查目录是否存在失败")
	}
	if !existDir {
		if err := fs.MkdirAll(filepath.Dir(filename), os.ModePerm); err != nil {
			return errors.Wrap(err, "创建目录失败")
		}
	}

	return nil
}

// GetFile 获取文档内容
func (repo Repo) GetFile(name string, filename string) ([]byte, error) {
	res, err := repo.workInWorktree(name, func(fs afero.Fs, conf *RepositoryConf) ([]byte, error) {
		return afero.ReadFile(fs, filepath.Join(conf.StorageDir, filename))
	})

	if err != nil {
		return nil, err
	}

	return res[0].([]byte), nil
}

// RepositoryTree 获取仓库目录树
func (repo Repo) RepositoryTree(name string, ext ...string) (*dtree.VueTree, error) {
	res, err := repo.workInWorktree(name, func(fs afero.Fs, conf *RepositoryConf) (*dtree.VueTree, error) {
		return dtree.BuildVueTreeForExt(fs, conf.StorageDir, ext...)
	})
	if err != nil {
		return nil, err
	}

	return res[0].(*dtree.VueTree), nil
}

func (repo Repo) workInRepository(name string, cb interface{}) ([]interface{}, error) {
	confProvider := func() (*RepositoryConf, error) {
		return repo.GetRepoConfByName(name)
	}

	repositoryProvider := func(ctx context.Context, fs billy.Filesystem, conf *RepositoryConf, writer *ProgressWriter) (*git.Repository, error) {
		return repo.openGitRepository(ctx, fs, conf, writer)
	}

	provider, err := repo.cc.ServiceProvider(confProvider, repositoryProvider)
	if err != nil {
		return nil, err
	}

	return repo.cc.CallWithProvider(cb, provider)
}

// OpenRepository 根据名称打开一个仓库
func (repo Repo) OpenRepository(name string) (*git.Repository, *RepositoryConf, error) {
	conf, err := repo.GetRepoConfByName(name)
	if err != nil {
		return nil, nil, errors.Wrap(err, "加载配置失败")
	}

	if conf == nil {
		return nil, nil, ErrRepositoryNotExist
	}

	var gitRepo *git.Repository
	if err := repo.cc.ResolveWithError(func(ctx context.Context, fs billy.Filesystem, writer *ProgressWriter) error {
		repository, err := repo.openGitRepository(ctx, fs, conf, writer)
		if err != nil {
			return err
		}

		gitRepo = repository
		return nil
	}); err != nil {
		return nil, conf, err
	}

	return gitRepo, conf, nil
}

// ResetRepository 重置仓库
func (repo Repo) ResetRepository(name string) error {
	conf, err := repo.GetRepoConfByName(name)
	if err != nil {
		return err
	}

	return repo.cc.ResolveWithError(func(fs billy.Filesystem) error {
		if conf.StorageDir == "" {
			return errors.New("仓库存储目录尚未设置，无法删除")
		}

		return fs.Remove(conf.StorageDir)
	})
}

// openGitRepository 打开一个 git 仓库
func (repo Repo) openGitRepository(ctx context.Context, fs billy.Filesystem, conf *RepositoryConf, writer *ProgressWriter) (*git.Repository, error) {
	worktreeFs, err := fs.Chroot(conf.StorageDir)
	if err != nil {
		return nil, errors.Wrapf(err, "无法打开目录 %s", conf.StorageDir)
	}

	gitFs, err := fs.Chroot(filepath.Join(conf.StorageDir, ".git"))
	if err != nil {
		return nil, errors.Wrapf(err, "无法打开目录 %s", filepath.Join(conf.StorageDir, ".git"))
	}

	stor := filesystem.NewStorage(gitFs, cache.NewObjectLRUDefault())

	gitRepo, err := git.Open(stor, worktreeFs)
	if err != nil {
		if err == git.ErrRepositoryNotExists {
			ctx, _ := context.WithTimeout(ctx, time.Second*60)
			gitRepo, err = git.CloneContext(ctx, stor, worktreeFs, &git.CloneOptions{
				URL:      conf.URL,
				Progress: writer,
				Auth:     conf.AuthConf.Auth(conf.Type),
			})
			if err != nil {
				return nil, errors.Wrap(err, "克隆 Git 仓库失败")
			}

			return gitRepo, nil
		}

		return nil, err
	}

	head, err := gitRepo.Head()
	if err != nil {
		return nil, err
	}

	log.Debugf("HEAD: %s", head.String())
	return gitRepo, nil
}

// Push 推送到远程仓库
func (repo *Repo) Push(name string) error {
	_, err := repo.workInRepository(name, func(ctx context.Context, repository *git.Repository, conf *RepositoryConf, writer *ProgressWriter) error {
		ctx, _ = context.WithTimeout(ctx, 30*time.Second)
		if err := repository.PushContext(ctx, &git.PushOptions{
			Auth:     conf.AuthConf.Auth(conf.Type),
			Progress: writer,
		}); err != nil {
			if err == git.NoErrAlreadyUpToDate {
				return nil
			}

			return errors.Wrap(err, "执行 git push 失败")
		}

		return nil
	})
	return err
}

// Pull 拉取更新
func (repo *Repo) Pull(name string) error {
	_, err := repo.workInRepository(name, func(ctx context.Context, repository *git.Repository, conf *RepositoryConf, writer *ProgressWriter) error {
		worktree, err := repository.Worktree()
		if err != nil {
			return errors.Wrap(err, "获取工作树失败")
		}

		if err := worktree.PullContext(ctx, &git.PullOptions{
			Auth:     conf.AuthConf.Auth(conf.Type),
			Progress: writer,
		}); err != nil {

			if err == git.NoErrAlreadyUpToDate {
				return nil
			}

			return errors.Wrap(err, "执行 git pull 失败")
		}

		return nil
	})
	return err
}

// commitChanges 提交文档变更
func (repo *Repo) commitChanges(name string, changed []string) error {
	_, err := repo.workInRepository(name, func(ctx context.Context, repository *git.Repository, conf *RepositoryConf) error {
		if len(changed) == 0 {
			return nil
		}

		worktree, err := repository.Worktree()
		if err != nil {
			return errors.Wrap(err, "文件已写入本地，但是获取工作树失败")
		}

		status, err := worktree.Status()
		if err != nil {
			return errors.Wrap(err, "文件已写入本地，但是获取仓库状态失败")
		}

		operations := make([]string, 0)
		for _, s := range array.StringUnique(changed) {
			s, _ = filepath.Rel(conf.StorageDir, s)
			fileStatus := status.File(s)

			// 有修改
			//     Worktree=Modified Staging=Unmodified
			// 新增文件
			//     Worktree=Untracked Staging=Untracked
			// 移动文件
			//     (old)Worktree=Deleted Staging=Unmodified
			//     (new)Worktree=Untracked Staging=Untracked

			opt := ""
			switch fileStatus.Worktree {
			case git.Modified:
				opt = "modified"
			case git.Untracked:
				opt = "added"
			case git.Deleted:
				opt = "deleted"
			case git.Renamed:
				opt = "renamed"
			default:
				opt = "unknown"
			}

			_, err := worktree.Add(s)
			if err != nil {
				return errors.Wrap(err, "文件已写入本地，但是添加到暂存区失败")
			}

			log.Debugf("[opt=%s] git add %s", opt, s)
			operations = append(operations, fmt.Sprintf("%s : %s", opt, s))
		}

		commitMessage := strings.Join(operations, ", ")
		commit, err := worktree.Commit(commitMessage, &git.CommitOptions{
			Author: &object.Signature{
				Name:  conf.AuthConf.Username,
				Email: conf.AuthConf.Email,
				When:  time.Now(),
			},
		})
		if err != nil {
			return errors.Wrap(err, "文件已写入本地，但是创建提交对象失败")
		}

		res, err := repository.CommitObject(commit)
		if err != nil {
			return errors.Wrap(err, "文件已写入本地，但是提交失败")
		}

		log.Debugf("[hash=%s] git commit -m %s", res.Hash.String(), commitMessage)

		return nil
	})

	return err
}

// buildRepositoryStorageDir 构建仓库存储目录名称
func (repo Repo) buildRepositoryStorageDir(repoURL string) (string, error) {
	rawURL, err := url.Parse(repoURL)
	if err != nil {
		return "", errors.Wrap(err, "parse repository url failed")
	}

	return strings.TrimLeft(strings.ReplaceAll(strings.ReplaceAll(rawURL.Path, "/", "-"), "\\", "-"), "-"), nil
}
