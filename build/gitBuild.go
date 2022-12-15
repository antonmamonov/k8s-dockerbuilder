package build

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type BuildConfig struct {
	DockerImageDestination string
}

func runCommand(name string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func Build(cfg *BuildConfig) error {
	// Check if we are in a git repository by checking if .git/config exists
	// read the .git/config file and get the repository name
	gitConfig, readGitConfigError := os.ReadFile("./.git/config")
	if readGitConfigError != nil {
		return readGitConfigError
	}

	// get the url value from the gitConfig file
	urlRegex := regexp.MustCompile(`url = (.*)`)
	urlMatches := urlRegex.FindStringSubmatch(string(gitConfig))

	if len(urlMatches) == 0 {
		return fmt.Errorf("Could not find the url in the git config file")
	}

	urlMatch := urlMatches[1]

	// get the branch name using 'git branch --show-current'
	branchName, branchNameError := runCommand("git", "branch", "--show-current")

	if branchNameError != nil {
		return branchNameError
	}

	// get the home directory path
	homeDir, homeDirError := os.UserHomeDir()

	if homeDirError != nil {
		return homeDirError
	}

	// get the docker config file from '~/.docker/config.json'
	dockerConfigBytes, readDockerConfigError := os.ReadFile(homeDir + "/.docker/config.json")

	if readDockerConfigError != nil {
		return readDockerConfigError
	}

	dockerConfigString := string(dockerConfigBytes)

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", homeDir+"/.kube/config")
	if err != nil {
		return err
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	clientCfg, _ := clientcmd.NewDefaultClientConfigLoadingRules().Load()
	namespace := clientCfg.Contexts[clientCfg.CurrentContext].Namespace

	if namespace == "" {
		namespace = "default"
	}

	jobName := cfg.DockerImageDestination

	// for kubernetes object names the valid regex for validation is regex used for validation is '[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*'
	// make sure the job name is valid and replace any invalid characters with an empty string
	jobNameRegex := regexp.MustCompile(`[^a-z0-9]`)
	jobName = "dockerbuild-" + jobNameRegex.ReplaceAllString(jobName, "")

	dockerBuildK8sJob := batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Job",
			APIVersion: "batch/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: jobName,
		},
		Spec: batchv1.JobSpec{
			BackoffLimit: func() *int32 { i := int32(7); return &i }(),
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					RestartPolicy: "OnFailure",
					Containers: []v1.Container{
						{
							Name:  "main",
							Image: "antonm/dockerbuilder",

							Command: []string{"/app/binutils/gitPullDockerBuildPush.sh"},

							// QUICK DEVELOPMENT
							// Command: []string{"/bin/bash"},
							// Args: []string{
							// 	"-c",
							// 	"sleep 9999;",
							// },

							Env: []v1.EnvVar{
								{
									Name:  "GIT_PULL_REPO_URL",
									Value: urlMatch,
								},
								{
									Name:  "GIT_BRANCH",
									Value: branchName,
								},
								{
									Name:  "DOCKER_IMAGE_DESTINATION",
									Value: cfg.DockerImageDestination,
								},
								{
									Name:  "DOCKER_CONFIG_JSON",
									Value: dockerConfigString,
								},
							},
						},
					},
				},
			},
		},
	}

	// create the syncK8sJob
	_, createSyncK8sJobError := clientset.BatchV1().Jobs(namespace).Create(context.TODO(), &dockerBuildK8sJob, metav1.CreateOptions{})

	if createSyncK8sJobError != nil {
		return createSyncK8sJobError
	}

	fmt.Println("[kubebuild] Successfully created the docker build job [" + jobName + "]")

	return nil
}
