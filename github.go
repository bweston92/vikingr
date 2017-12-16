package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/go-github/github"
)

type repositoryJobError map[*github.Repository]error

func (e repositoryJobError) Error() string {
	l := make([]string, len(e))
	for r, err := range e {
		l = append(l, fmt.Sprintf("repo [%s] err: %s", r.GetFullName(), err.Error()))
	}

	return strings.Join(l, "\n")
}

func runCheck(ctx context.Context, gh *github.Client, org string) error {
	rs, err := getRepositories(ctx, gh, org)
	if err != nil {
		return err
	}

	errs := repositoryJobError{}
	for _, r := range rs {
		debug("running check on repository %s", r.GetName())
		if err := checkRepository(ctx, gh, r); err != nil {
			errs[r] = err
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func getGithubClient(u, t string) *github.Client {
	return github.NewClient((&github.BasicAuthTransport{
		Username: u,
		Password: t,
	}).Client())
}

func getRepositories(ctx context.Context, gh *github.Client, org string) ([]*github.Repository, error) {
	opts := &github.RepositoryListByOrgOptions{}
	opts.PerPage = 250
	rs, _, err := gh.Repositories.ListByOrg(ctx, org, opts)

	if err != nil {
		debug("issue fetching repositories %s", err)
	} else {
		debug("fetched %d repositories for %s organsiation", len(rs), org)
	}

	return rs, err
}

func checkRepository(ctx context.Context, gh *github.Client, r *github.Repository) error {
	if !requiresProtection(r) {
		debug("no check required")
		return nil
	}

	b := r.GetDefaultBranch()
	if len(b) == 0 {
		return errors.New("no default branch")
	}
	debug("repository default branch: %s", b)

	p, res, err := gh.Repositories.GetBranchProtection(ctx, r.Owner.GetLogin(), r.GetName(), b)

	if err != nil {
		if res.StatusCode != http.StatusNotFound {
			debug("unknown error getting branch permissions: %s", err)
			return err
		}

		debug("branch has no permissions, attempting to protect %s branch", b)
		return protectBranch(ctx, gh, r, b)
	}

	if !isValidProtection(p) {
		debug("branch permissions are not valid, just going to tinker them")
		return fixBranch(ctx, gh, r, b, p)
	}

	debug("branch permissions perfectly fine")
	return nil
}

func requiresProtection(r *github.Repository) bool {
	return r.GetPrivate()
}

func isValidProtection(p *github.Protection) bool {
	return p.EnforceAdmins.Enabled && p.RequiredPullRequestReviews != nil && p.RequiredPullRequestReviews.DismissStaleReviews
}

func protectBranch(ctx context.Context, gh *github.Client, r *github.Repository, b string) error {
	_, _, err := gh.Repositories.UpdateBranchProtection(ctx, r.Owner.GetLogin(), r.GetName(), b, &github.ProtectionRequest{
		RequiredPullRequestReviews: &github.PullRequestReviewsEnforcementRequest{
			DismissStaleReviews: true,
			DismissalRestrictionsRequest: &github.DismissalRestrictionsRequest{
				Users: []string{},
				Teams: []string{},
			},
		},
		EnforceAdmins: true,
	})
	return err
}

func fixBranch(ctx context.Context, gh *github.Client, r *github.Repository, b string, p *github.Protection) error {
	req := &github.ProtectionRequest{
		RequiredStatusChecks: p.RequiredStatusChecks,
		RequiredPullRequestReviews: &github.PullRequestReviewsEnforcementRequest{
			DismissStaleReviews:          true,
			DismissalRestrictionsRequest: getDismissalRestrictionsRequestForExistingProtection(p),
		},
		EnforceAdmins: true,
	}

	_, _, err := gh.Repositories.UpdateBranchProtection(ctx, r.Owner.GetLogin(), r.GetName(), b, req)
	return err
}

func getDismissalRestrictionsRequestForExistingProtection(p *github.Protection) *github.DismissalRestrictionsRequest {

	if p.Restrictions != nil {
		u := make([]string, len(p.Restrictions.Users))
		t := make([]string, len(p.Restrictions.Teams))

		for i, user := range p.Restrictions.Users {
			u[i] = user.GetLogin()
		}

		for i, team := range p.Restrictions.Teams {
			t[i] = team.GetName()
		}

		return &github.DismissalRestrictionsRequest{
			Users: u,
			Teams: t,
		}
	}

	return &github.DismissalRestrictionsRequest{
		Users: []string{},
		Teams: []string{},
	}
}
