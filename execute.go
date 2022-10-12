package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/adnsv/go-utils/git"
	"github.com/adnsv/go-utils/prompt"
	"github.com/adnsv/go-utils/version"
	cli "github.com/jawher/mow.cli"
)

type exec_options struct {
	// prefix for new tags
	prefix string

	// allow tagging of a repository that contains uncommited changes
	allow_dirty bool

	// show quad-related
	with_quad bool
}

func (opts *exec_options) bind_cli(cli *cli.Cli) {
	cli.StringOptPtr(&opts.prefix, "p prefix", "AUTO", "prefix for new tags")
	cli.BoolOptPtr(&opts.allow_dirty, "d allow-dirty", false, "allow tagging of repos that contain uncommited changes")
}

var errDirtyRepo = errors.New("the repository has uncommited changes")

func execute(opts *exec_options) error {
	wd, stats, err := get_stats()

	if err == git.ErrNoTags {
		return create_first_tag(opts)
	} else if err != nil {
		return fmt.Errorf("failed to obtain git stats: %w", err)
	}

	oldtag := stats.Description.Tag
	vi, err := git.ParseVersion(stats.Description)
	if stats.Dirty {
		print_keyval("last tag", oldtag)
		if stats.Description.AdditionalCommits > 0 {
			print_keyval("additional commits", stats.Description.AdditionalCommits)
		}

		if !opts.allow_dirty {
			fmt.Println()
			fmt.Println(fmt_bold("WARNING") + ": modified since the last commit")
			fmt.Println("         (execute 'git status' for mode detail)")
			fmt.Println()
			fmt.Println("Commit your changes before updating the tag")
			print_dim("or re-run with --allow-dirty to force tagging in dirty state")
			fmt.Println()
			return errDirtyRepo
		} else {
			print_keyval("state", fmt_bold("dirty")+", has uncommited changes")
		}
	} else {
		print_keyval("state", "clean, no uncommited changes")
	}

	if !opts.with_quad && err == version.ErrNumberOfAdditionalCommitsIsTooLarge {
		err = nil
	}

	if err != nil {
		vi, oldtag, err = handle_tag_parse_error(wd, vi, oldtag, err)
		if err != nil {
			return err
		}
		print_keyval("last semantic tag", oldtag)
	}

	if opts.prefix == "AUTO" {
		opts.prefix = "v"
		if len(oldtag) > 1 {
			d := strings.IndexAny(oldtag, "0123456789")
			if d == 0 {
				opts.prefix = ""
			} else if d == 1 && oldtag[0] == 'v' || oldtag[0] == 'V' {
				opts.prefix = oldtag[:1]
			}
		}

		if opts.prefix == "" {
			print_keyval("auto prefix", "no")
		} else {
			print_keyval("auto prefix", opts.prefix)
		}
	}

	stats.Description.Tag = strings.TrimPrefix(stats.Description.Tag, opts.prefix)

	if stats.Description.AdditionalCommits > 0 {
		print_keyval("additional commits", stats.Description.AdditionalCommits)
	}

	print_keyval("semantic ver", vi.Semantic)
	if opts.with_quad {
		print_keyval("version quad", vi.Quad.String())
	}

	if stats.Description.AdditionalCommits == 0 {
		fmt.Println()
		fmt.Printf("The current state of repository is already tagged as %s.\n", fmt_tag(oldtag))
		fmt.Println("If you proceed, you will have more than one tag pointing to the same state.")
		fmt.Println()
		if !prompt.YN("Still want to proceed [y/n]?") {
			return errUserCancelled
		}
	}

	actions := collectActions(vi.Semantic)
	if len(actions) == 0 {
		fmt.Println("No actions available")
		return nil
	}

	fmt.Println()

	choices := make([]string, 0, len(actions))
	for _, a := range actions {
		desc := a.desc
		comment := ""
		if i := strings.IndexByte(desc, '|'); i >= 0 {
			comment = " " + fmt_dim("("+desc[i+1:]+")")
			desc = desc[:i]
		}
		if a.showPRchoice {
			choices = append(choices, fmt.Sprintf("%s %s%s ...",
				desc, fmt_tag(opts.prefix+a.ver.String()), comment))
		} else {
			choices = append(choices, fmt.Sprintf("%s %s%s",
				desc, fmt_tag(opts.prefix+a.ver.String()), comment))
		}
	}

	choice := prompt.Choose("available actions:", choices...)
	action := actions[choice-1]
	newver := action.ver

	if action.showPRchoice {
		fmt.Println()
		fmt.Println("Select (pre-)release type")
		fmt.Printf("- 'alpha'   for %s\n", fmt_tag(opts.prefix+withPR(action.ver, "alpha", 1).String()))
		fmt.Printf("- 'beta'    for %s\n", fmt_tag(opts.prefix+withPR(action.ver, "beta", 1).String()))
		fmt.Printf("- 'rc'      for %s\n", fmt_tag(opts.prefix+withPR(action.ver, "rc", 1).String()))
		fmt.Printf("- 'release' for %s\n", fmt_tag(opts.prefix+withoutPR(action.ver).String()))

		choice := prompt.Enum("type", "alpha", "beta", "rc", "release")
		if choice != "release" {
			newver.Pre = makePR(choice, 1)
		} else {
			newver.Pre = newver.Pre[:0]
		}
	}

	tag := opts.prefix + newver.String()
	comment := generate_tag_comment(tag)
	return perform_tagging(tag, comment)
}

func generate_tag_comment(tag string) string {
	return fmt.Sprintf("tagging as %s", tag)
}

func create_first_tag(opts *exec_options) error {
	if opts.prefix == "AUTO" {
		opts.prefix = "v"
	}
	tag := opts.prefix + "0.0.1"
	comment := generate_tag_comment(tag)

	fmt.Println()
	fmt.Println(fmt_bold("WARNING") + ": no existing tags found")
	fmt.Println()
	fmt.Println("Proceed to create the first tag with this utility,")
	fmt.Println("or cancel and assign it manually:")
	fmt.Println()
	fmt.Printf("    git tag -a %s -m \"%s\"\n", tag, comment)
	fmt.Println()

	return perform_tagging(tag, comment)
}

func run_git_tag(tag, comment string) error {
	begin_dim()
	defer end_dim()

	fmt.Printf("executing: 'git tag -a %s -m \"%s\"\n", tag, comment)
	cmd := exec.Command("git", "tag", "-a", tag, "-m", comment)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func run_git_push_origin(tag string) error {
	begin_dim()
	defer end_dim()

	fmt.Printf("executing: 'git push origin %s'\n", tag)
	cmd := exec.Command("git", "push", "origin", tag)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func perform_tagging(tag, comment string) error {

	tag_comment := fmt.Sprintf("tagging as %s", tag)

	fmt.Println()
	fmt.Printf("Ready to tag as %s %s:\n", fmt_tag(tag), fmt_dim("(with comment '"+tag_comment+"')"))
	fmt.Println()
	if !prompt.YN("proceed [y/n]? ") {
		return errUserCancelled
	}

	fmt.Println()
	err := run_git_tag(tag, comment)
	if err != nil {
		return fmt.Errorf("failed to execute git tag: %w", err)
	}

	fmt.Println()
	fmt.Printf("Your local repository is now tagged as %s\n", fmt_tag(tag))
	fmt.Println()
	fmt.Printf("To push this change to remote, execute:\n\n")
	fmt.Printf("    git push origin %s\n\n", tag)
	begin_dim()
	fmt.Println("NOTE: to revert local and/or remote tagging,")
	fmt.Println("      re-run rtag with --undo")
	end_dim()
	fmt.Println()
	fmt.Println("This utility can push the new tag for you")
	fmt.Println()
	if !prompt.YN("Proceed with push [y/n]? ") {
		return errUserCancelled
	}

	fmt.Println()
	err = run_git_push_origin(tag)
	if err != nil {
		return fmt.Errorf("failed to execute git push: %w", err)
	}
	return err
}

var errNoSemanticTagsFound = errors.New("no semantic tags found")

func handle_tag_parse_error(workdir string, last_vi *git.VersionInfo, last_tag string, parse_err error) (*git.VersionInfo, string, error) {
	fmt.Printf("tag parse error: %s\n", parse_err)
	if parse_err == version.ErrNumberOfAdditionalCommitsIsTooLarge {

		fmt.Printf(fmt_bold("WARNING") + `: generation of version quads fails if the number of
additional commits exceeds 99.

Please consider bumping up the version to resolve this issue.

If you choose to proceed, it will be capped to 99 in the generated quad.
`)
		if prompt.YN("Do you want to proceed [y/n]?") {
			return last_vi, last_tag, nil
		} else {
			return last_vi, last_tag, errUserCancelled
		}
	}

	sem_tag, sem_vi, err := git.LastSemanticTag(workdir)
	if err == nil {
		fmt.Printf(fmt_bold("WARNING")+": last tag '%s' does not conform to semantic version syntax\n", last_tag)
		fmt.Printf("however, there is an older tag '%s' that can be used instead\n", sem_tag)
		fmt.Printf("\n")
		if prompt.YN("Proceed with '" + sem_tag + "' as base [y/n]?") {
			return sem_vi, sem_tag, nil
		} else {
			return sem_vi, sem_tag, errUserCancelled
		}
	}
	fmt.Println()
	fmt.Printf(fmt_bold("ERROR")+": last tag '%s' does not conform to semantic version syntax\n", last_tag)
	fmt.Println()
	fmt.Printf("This utility expects your repository to be tagged with semantic tags,\n")
	fmt.Printf("see https://semver.org for more information\n")
	fmt.Println()
	fmt.Printf("Exiting now\n")
	return last_vi, last_tag, errNoSemanticTagsFound
}
