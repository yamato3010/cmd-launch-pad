package git

import (
	"fmt"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"time"
)

// GitManager はGitリポジトリ操作を提供するラッパー
type GitManager struct {
	repoPath string
	repo     *gogit.Repository
}

// NewGitManager は指定パスのGitManagerを生成する。リポジトリが存在しない場合はエラー。
func NewGitManager(repoPath string) (*GitManager, error) {
	repo, err := gogit.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("Gitリポジトリのオープンに失敗しました: %w", err)
	}
	return &GitManager{repoPath: repoPath, repo: repo}, nil
}

// Init は指定パスをGitリポジトリとして初期化する
func Init(repoPath string) (*GitManager, error) {
	repo, err := gogit.PlainInit(repoPath, false)
	if err != nil {
		// 既にGitリポジトリの場合はオープンを試みる
		repo, err = gogit.PlainOpen(repoPath)
		if err != nil {
			return nil, fmt.Errorf("Gitリポジトリの初期化に失敗しました: %w", err)
		}
	}
	return &GitManager{repoPath: repoPath, repo: repo}, nil
}

// SetRemote はリモートリポジトリURLを設定する
func (g *GitManager) SetRemote(name, url string) error {
	// 既存のリモートを削除してから再作成
	_ = g.repo.DeleteRemote(name)
	_, err := g.repo.CreateRemote(&config.RemoteConfig{
		Name: name,
		URLs: []string{url},
	})
	if err != nil {
		return fmt.Errorf("リモートの設定に失敗しました: %w", err)
	}
	return nil
}

// AddAll はワーキングツリーの全変更をステージングする
func (g *GitManager) AddAll() error {
	wt, err := g.repo.Worktree()
	if err != nil {
		return fmt.Errorf("ワーキングツリーの取得に失敗しました: %w", err)
	}
	if err := wt.AddWithOptions(&gogit.AddOptions{All: true}); err != nil {
		return fmt.Errorf("git addに失敗しました: %w", err)
	}
	return nil
}

// Commit は変更をコミットする
func (g *GitManager) Commit(message string) error {
	wt, err := g.repo.Worktree()
	if err != nil {
		return fmt.Errorf("ワーキングツリーの取得に失敗しました: %w", err)
	}
	_, err = wt.Commit(message, &gogit.CommitOptions{
		Author: &object.Signature{
			Name:  "cmd-launch-pad",
			Email: "cmd-launch-pad@local",
			When:  time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("コミットに失敗しました: %w", err)
	}
	return nil
}

// Push はリモートにプッシュする
func (g *GitManager) Push(remote, branch string, auth *http.BasicAuth) error {
	opts := &gogit.PushOptions{
		RemoteName: remote,
	}
	if auth != nil {
		opts.Auth = auth
	}
	if err := g.repo.Push(opts); err != nil && err != gogit.NoErrAlreadyUpToDate {
		return fmt.Errorf("プッシュに失敗しました: %w", err)
	}
	return nil
}

// Pull はリモートからプルする
func (g *GitManager) Pull(remote, branch string, auth *http.BasicAuth) error {
	wt, err := g.repo.Worktree()
	if err != nil {
		return fmt.Errorf("ワーキングツリーの取得に失敗しました: %w", err)
	}
	opts := &gogit.PullOptions{
		RemoteName: remote,
	}
	if auth != nil {
		opts.Auth = auth
	}
	if err := wt.Pull(opts); err != nil && err != gogit.NoErrAlreadyUpToDate {
		return fmt.Errorf("プルに失敗しました: %w", err)
	}
	return nil
}

// Status はワーキングツリーの状態を文字列で返す
func (g *GitManager) Status() (string, error) {
	wt, err := g.repo.Worktree()
	if err != nil {
		return "", fmt.Errorf("ワーキングツリーの取得に失敗しました: %w", err)
	}
	status, err := wt.Status()
	if err != nil {
		return "", fmt.Errorf("ステータスの取得に失敗しました: %w", err)
	}
	return status.String(), nil
}

// IsClean はワーキングツリーがクリーンかどうかを返す
func (g *GitManager) IsClean() (bool, error) {
	wt, err := g.repo.Worktree()
	if err != nil {
		return false, err
	}
	status, err := wt.Status()
	if err != nil {
		return false, err
	}
	return status.IsClean(), nil
}
