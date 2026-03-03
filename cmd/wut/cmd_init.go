package main

import (
	"fmt"
	"os"

	"github.com/simonbs/wut/src/context"
	"github.com/simonbs/wut/src/git"
	"github.com/simonbs/wut/src/worktree"
)

func cmdInit() {
	binPath, _ := os.Executable()
	wrapper := `wut() {
  local wut_bin="` + binPath + `"
  export WUT_WRAPPER_ACTIVE=1
  local output
  output=$("$wut_bin" "$@" 2>&1)
  local exit_code=$?
  local cd_marker
  cd_marker=$(echo "$output" | grep "^__WUT_CD__:" | head -1)
  if [ -n "$cd_marker" ]; then
    local target_dir="${cd_marker#__WUT_CD__:}"
    if [ -d "$target_dir" ]; then
      cd "$target_dir" || return 1
    fi
    local filtered
    filtered=$(printf "%s" "$output" | grep -v "^__WUT_CD__:")
    if [[ -n "${filtered//[[:space:]]/}" ]]; then
      printf "%s\\n" "$filtered"
    fi
  else
    if [[ -n "${output//[[:space:]]/}" ]]; then
      printf "%s\\n" "$output"
    fi
  fi
  return $exit_code
}

_wut_completions() {
  local cur="${COMP_WORDS[COMP_CWORD]}"
  local prev="${COMP_WORDS[COMP_CWORD-1]}"
  local cmd="${COMP_WORDS[1]}"

  _wut_complete_worktree_targets() {
    if [[ "$cur" == /* || "$cur" == .* || "$cur" == "~"* ]]; then
      local paths
      paths=$("` + binPath + `" --completions paths 2>/dev/null)
      COMPREPLY=($(compgen -W "$paths" -- "$cur"))
      return
    fi

    local branches
    branches=$("` + binPath + `" --completions branches 2>/dev/null)
    COMPREPLY=($(compgen -W "$branches" -- "$cur"))
  }

  _wut_complete_removable_targets() {
    if [[ "$cur" == /* || "$cur" == .* || "$cur" == "~"* ]]; then
      local paths
      paths=$("` + binPath + `" --completions removable-paths 2>/dev/null)
      COMPREPLY=($(compgen -W "$paths" -- "$cur"))
      return
    fi

    local branches
    branches=$("` + binPath + `" --completions removable-branches 2>/dev/null)
    COMPREPLY=($(compgen -W "$branches" -- "$cur"))
  }
  
  if [[ ${COMP_CWORD} -eq 1 ]]; then
    COMPREPLY=($(compgen -W "new mv list go path rm" -- "$cur"))
    return
  fi

  if [[ "$cmd" == "new" ]]; then
    if [[ "$prev" == "--from" ]]; then
      local refs
      refs=$("` + binPath + `" --completions refs 2>/dev/null)
      COMPREPLY=($(compgen -W "$refs" -- "$cur"))
      return
    fi

    if [[ "$cur" == -* ]]; then
      COMPREPLY=($(compgen -W "--from" -- "$cur"))
      return
    fi
  fi
  
  case "$prev" in
    mv|go|path|rm)
      if [[ "$prev" == "go" ]]; then
        _wut_complete_worktree_targets
      elif [[ "$prev" == "rm" ]]; then
        _wut_complete_removable_targets
      else
        local branches
        branches=$("` + binPath + `" --completions branches 2>/dev/null)
        COMPREPLY=($(compgen -W "$branches" -- "$cur"))
      fi
      return
  esac

  if [[ "$cmd" == "rm" && "$cur" == -* ]]; then
    COMPREPLY=($(compgen -W "--force" -- "$cur"))
    return
  fi

  if [[ "$cmd" == "go" || "$cmd" == "rm" ]]; then
    if [[ "$cmd" == "go" ]]; then
      _wut_complete_worktree_targets
    else
      _wut_complete_removable_targets
    fi
    return
  fi
}

complete -F _wut_completions wut

# zsh completion
if [[ -n ${ZSH_VERSION-} ]]; then
  autoload -U +X bashcompinit && bashcompinit
fi`
	fmt.Println(wrapper)
}

func cmdCompletions(args []string) {
	if len(args) < 1 {
		return
	}

	switch args[0] {
	case "branches":
		ctx, err := context.Create()
		if err != nil {
			return
		}
		entries, _ := worktree.ParseList(ctx.RepoRoot)
		for _, e := range entries {
			if e.BranchName != "" {
				fmt.Println(e.BranchName)
			}
		}
	case "paths":
		ctx, err := context.Create()
		if err != nil {
			return
		}
		entries, err := worktree.ParseList(ctx.RepoRoot)
		if err != nil {
			return
		}
		for _, e := range entries {
			if e.Path != "" {
				fmt.Println(e.Path)
			}
		}
	case "removable-branches":
		ctx, err := context.Create()
		if err != nil {
			return
		}
		entries, err := worktree.ParseList(ctx.RepoRoot)
		if err != nil {
			return
		}
		for _, e := range entries {
			if e.Path == ctx.RepoRoot {
				continue
			}
			if e.BranchName != "" {
				fmt.Println(e.BranchName)
			}
		}
	case "removable-paths":
		ctx, err := context.Create()
		if err != nil {
			return
		}
		entries, err := worktree.ParseList(ctx.RepoRoot)
		if err != nil {
			return
		}
		for _, e := range entries {
			if e.Path == "" || e.Path == ctx.RepoRoot {
				continue
			}
			fmt.Println(e.Path)
		}
	case "refs":
		ctx, err := context.Create()
		if err != nil {
			return
		}
		refs, err := git.ListBranchRefs(ctx.RepoRoot)
		if err != nil {
			return
		}
		for _, ref := range refs {
			fmt.Println(ref)
		}
	}
}
