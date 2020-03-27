package repo

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/C123R/helm-blob/pkg/blob"
	"gopkg.in/yaml.v2"
	"helm.sh/helm/pkg/chart/loader"
	"helm.sh/helm/pkg/helmpath"
	"helm.sh/helm/pkg/provenance"
	helmrepo "helm.sh/helm/pkg/repo"
)

type IndexFile helmrepo.IndexFile

type Repo struct {
	name, url string
	blob      blob.BlobConnect
}

const index = "index.yaml"

// NewRepoByName returns Repo with blobConnect by repoName
func NewRepoByName(repoName string) (Repo, error) {

	repoUrl, err := getRepoUrl(repoName)
	if err != nil {
		return Repo{}, err
	}
	b, err := blob.NewBlobConnect(repoUrl)
	if err != nil {
		return Repo{}, fmt.Errorf("Unable to connect to the Blob %s: %v", repoUrl, err)
	}
	return Repo{
		name: repoName,
		url:  repoUrl,
		blob: b,
	}, nil
}

// NewRepoByUrl returns Repo with blobConnect by repoUrl
func NewRepoByUrl(repoUrl string) (Repo, error) {

	b, err := blob.NewBlobConnect(repoUrl)
	if err != nil {
		return Repo{}, fmt.Errorf("Unable to connect to the Blob %s: %v", repoUrl, err)
	}
	return Repo{
		url:  repoUrl,
		blob: b,
	}, nil
}

// Init initialize new repo by adding default index.yaml
func (r Repo) Init() error {

	// Check if specified repository is valid ( by checking index.yaml)
	if r.blob.IndexFileExists() {
		return fmt.Errorf("Specified %s chart repository is already initialized!", r.url)
	}
	indexFile := helmrepo.NewIndexFile()
	indexInByte, err := yaml.Marshal(indexFile)
	if err != nil {
		return err
	}
	if err := r.blob.Upload(index, indexInByte); err != nil {
		return err
	}
	fmt.Printf(`Successfully initialized %s.

Now you can add this repo using 'helm repo add <repo-name> %s'`, r.url, r.url)
	fmt.Println()
	return nil
}

// Push uploads provided chart, updates index.yaml
// with new version of index.yaml on remote blob repo
func (r Repo) Push(chartpath string, force bool) error {

	repoIndex, err := r.fetchIndexFile()
	if err != nil {
		return fmt.Errorf("Looks like %s is not a valid chart repository or cannot be reached: %v", r.name, err)
	}
	c, err := loader.Load(chartpath)
	if err != nil {
		return fmt.Errorf("Error loading chart: %v", err)
	}

	// check if specified chart with version is present or not?
	if repoIndex.Has(c.Metadata.Name, c.Metadata.Version) && !force {
		return fmt.Errorf("%s with version %s is already present. Use --force/-f for force uploading!", c.Metadata.Name, c.Metadata.Version)
	} else if force {
		chartVersions := repoIndex.Entries[c.Metadata.Name]
		if len(chartVersions) != 0 {
			for i, v := range chartVersions {
				if c.Metadata.Version == v.Version {
					chartVersions = append(chartVersions[:i], chartVersions[i+1:]...)
					break
				}
			}
			repoIndex.Entries[c.Metadata.Name] = chartVersions
		}
	}

	// Generate SHA256 for the provided chart
	hash, err := provenance.DigestFile(chartpath)
	if err != nil {
		return err
	}
	repoIndex.Add(c.Metadata, getChartFileName(chartpath), r.url, hash)
	if err := r.updateIndex(repoIndex); err != nil {
		return err
	}
	if err := r.uploadChart(chartpath); err != nil {
		return err
	}
	fmt.Printf("Successfully uploaded chart %s in %s\n", getChartFileName(chartpath), r.name)
	return nil
}

// Delete deletes provided chart, updates index.yaml
// by deleting entry for specfic chart on remote blob repo
func (r Repo) Delete(chartName string, version string) error {

	var charts2delete []string
	repoIndex, err := r.fetchIndexFile()
	if err != nil {
		return fmt.Errorf("Looks like %s is not a valid chart repository or cannot be reached: %v", r.name, err)
	}
	// check if specified chart with version is present or not?
	chartVersions, ok := repoIndex.Entries[chartName]
	if !ok {
		return fmt.Errorf("No chart found with name %s\n", chartName)
	}

	if len(chartVersions) != 0 {
		for i, v := range chartVersions {
			if version == v.Version {
				chartVersions = append(chartVersions[:i], chartVersions[i+1:]...)
			}
			if version == v.Version || version == "" {
				charts2delete = append(charts2delete, fmt.Sprintf("%s-%s.tgz", chartName, v.Version))
			}
		}
		repoIndex.Entries[chartName] = chartVersions
	}

	if len(chartVersions) == 0 || version == "" {
		delete(repoIndex.Entries, chartName)
	}

	if len(charts2delete) == 0 {
		return fmt.Errorf("No chart found with name %s and version %s\n", chartName, version)
	}

	for _, chart := range charts2delete {
		if err := r.blob.Delete(chart); err != nil {
			return err
		}
		fmt.Printf("Successfully deleted chart %s \n", chart)
	}

	if err := r.updateIndex(repoIndex); err != nil {
		return err
	}
	return nil
}

// Fetch gets file for provided repoURL
func (r Repo) Fetch() error {

	repoIndex, err := r.blob.Download(index)
	if err != nil {
		return fmt.Errorf("Looks like %s is not a valid chart repository or cannot be reached: %v", r.url, err)
	}
	fmt.Println(repoIndex.String())
	return nil
}

// fetchIndexFile downloads "index.yaml" from remote repository(blob)
// returns sorted index file
func (r Repo) fetchIndexFile() (*helmrepo.IndexFile, error) {

	indexFile := helmrepo.NewIndexFile()
	repoIndex, err := r.blob.Download(index)
	if err != nil {
		return indexFile, err
	}
	err = yaml.Unmarshal(repoIndex.Bytes(), indexFile)
	if err != nil {
		return indexFile, err
	}
	indexFile.SortEntries()
	return indexFile, nil
}

func (r Repo) updateIndex(indexfile *helmrepo.IndexFile) error {

	indexfileInByte, err := yaml.Marshal(indexfile)
	if err != nil {
		return err
	}
	if err := r.blob.Upload(index, indexfileInByte); err != nil {
		return err
	}
	return nil
}

func (r Repo) uploadChart(chartpath string) error {

	data, err := ioutil.ReadFile(chartpath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("Provided chart %s doesnot exist: %v", chartpath, err)
		}
		return fmt.Errorf("Failed to read %s: %v", chartpath, err)
	}

	if err := r.blob.Upload(getChartFileName(chartpath), data); err != nil {
		return err
	}
	return nil
}

func getRepoUrl(repoName string) (string, error) {

	repositoryConfigPath := envOr("HELM_REPOSITORY_CONFIG", helmpath.ConfigPath("repositories.yaml"))
	var repoUrl string
	f, err := helmrepo.LoadFile(repositoryConfigPath)
	if err != nil {
		return repoUrl, err
	}
	for _, r := range f.Repositories {
		if r.Name == repoName {
			return r.URL, nil
		}
	}
	return repoUrl, fmt.Errorf("couldn't find repository specified repository")
}
