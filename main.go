package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	contextCMD "github.com/okteto/okteto/cmd/context"
	"github.com/okteto/okteto/cmd/utils"
	"github.com/okteto/okteto/pkg/log"
	"github.com/okteto/okteto/pkg/model"
	"github.com/okteto/okteto/pkg/okteto"
)

func main() {
	ctx := context.Background()

	if err := initOktetoContext(ctx); err != nil {
		log.Fail(err.Error())
		os.Exit(1)
	}

	previewName := ""
	if len(os.Args) > 1 {
		previewName = os.Args[1]
	} else {
		log.Fail("Preview name is needed")
		os.Exit(1)
	}

	oktetoClient, err := okteto.NewOktetoClient()
	if err != nil {
		log.Fail(err.Error())
		os.Exit(1)
	}

	_, err = deployPreview(ctx, previewName, oktetoClient)
	if err != nil {
		log.Fail(err.Error())
		os.Exit(1)
	}

	log.Information("Preview URL: %s", getPreviewURL(previewName))

}

func initOktetoContext(ctx context.Context) error {
	oktetoToken := os.Getenv("OKTETO_TOKEN")
	oktetoURL := os.Getenv("OKTETO_URL")

	ctxOptions := &contextCMD.ContextOptions{
		Token:   oktetoToken,
		Context: oktetoURL,
	}
	if err := contextCMD.Run(ctx, ctxOptions); err != nil {
		return err
	}
	return nil
}

func deployPreview(ctx context.Context, name string, oktetoClient *okteto.OktetoClient) (*okteto.PreviewResponse, error) {
	repository := os.Getenv("REPOSITORY")
	branch := os.Getenv("BRANCH")
	scope := os.Getenv("SCOPE")
	sourceURL := os.Getenv("SOURCE_URL")
	filename := os.Getenv("FILENAME")
	vars := os.Getenv("VARIABLES")

	repository, err := getRepository(ctx, repository)
	if err != nil {
		return nil, err
	}
	branch, err = getBranch(ctx, branch)
	if err != nil {
		return nil, err
	}

	varList := []okteto.Variable{}
	if len(vars) > 0 {
		variables := strings.Split(vars, ";")
		for _, v := range variables {
			kv := strings.SplitN(v, "=", 2)
			if len(kv) != 2 {
				return nil, fmt.Errorf("invalid variable value '%s': must follow KEY=VALUE format", v)
			}
			varList = append(varList, okteto.Variable{
				Name:  kv[0],
				Value: kv[1],
			})
		}
	}

	fmt.Printf("  Name: %s\n  Repo: %s @ %s\n  Scope: %s\n  Variables: %s\n", name, repository, branch, scope, vars)

	resp, err := oktetoClient.DeployPreview(ctx, name, scope, repository, branch, sourceURL, filename, varList)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func getRepository(ctx context.Context, repository string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get the current working directory: %w", err)
	}

	if repository == "" {
		r, err := model.GetRepositoryURL(cwd)

		if err != nil {
			return "", err
		}

		repository = r
	}
	return repository, nil
}

func getBranch(ctx context.Context, branch string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get the current working directory: %w", err)
	}
	if branch == "" {
		b, err := utils.GetBranch(ctx, cwd)

		if err != nil {
			return "", err
		}
		branch = b
	}
	return branch, nil
}

func getPreviewURL(name string) string {
	oktetoURL := okteto.Context().Name
	previewURL := fmt.Sprintf("%s/#/previews/%s", oktetoURL, name)
	return previewURL
}
