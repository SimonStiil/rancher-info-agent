package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	//	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/rancher/norman/clientbase"
	"github.com/rancher/norman/types"
	v3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	management "github.com/rancher/rancher/pkg/client/generated/management/v3"
	//	"k8s.io/apimachinery/pkg/runtime/serializer"
)

// https://github.com/rancher/rancher/blob/00ffe87e7cf71881226c798eb3fb8bfece4ceb4b/tests/framework/clients/rancher/client.go#L75
func newClientFor(restConfig *rest.Config) (*management.Client, error) {
	return management.NewClient(clientOpts(restConfig))
}

func clientOpts(restConfig *rest.Config) *clientbase.ClientOpts {
	return &clientbase.ClientOpts{
		URL:        fmt.Sprintf("https://%s/apis", restConfig.Host),
		TokenKey:   restConfig.BearerToken,
		Insecure:   restConfig.Insecure,
		HTTPClient: &http.Client{Transport: restConfig.Transport},
		CACerts:    restConfig.TLSClientConfig.CAFile,
	}
}

/*
	type CattleManagementV3Interface interface {
		// https://github.com/rancher/rancher/blob/release/v2.8/pkg/apis/management.cattle.io/v3/authz_types.go#L31
		Project(namespace string) ProjectInterface
		// https://github.com/rancher/rancher/blob/release/v2.8/pkg/apis/management.cattle.io/v3/cluster_types.go#L100C6-L100C13
		Cluster() ClusterInterface
	}

	type CattleManagementV3Client struct {
		restClient rest.Interface
	}

// https://www.martin-helmich.de/en/blog/kubernetes-crd-client.html
// https://devpress.csdn.net/k8s/62ebdfa589d9027116a0fa54.html

	func NewForConfig(c *rest.Config) (CattleManagementV3Interface, error) {
		config := *c
		config.ContentConfig.GroupVersion = &v3.SchemeGroupVersion
		config.APIPath = "/apis"
		config.NegotiatedSerializer = serializer.NewCodecFactory(scheme.Scheme)
		config.UserAgent = rest.DefaultKubernetesUserAgent()

		client, err := rest.RESTClientFor(&config)
		if err != nil {
			return nil, err
		}
		return &CattleManagementV3Client{restClient: client}, nil
	}

	func (c *CattleManagementV3Client) Project(namespace string) ProjectInterface {
		return &ProjectClient{
			restClient: c.restClient,
			ns:         namespace,
		}
	}

	type ProjectInterface interface {
		List(ctx context.Context, opts metav1.ListOptions) (*v3.ProjectList, error)
		Get(ctx context.Context, name string, options metav1.GetOptions) (*v3.Project, error)
		Create(ctx context.Context, project *v3.Project) (*v3.Project, error)
		Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	}

	type ProjectClient struct {
		restClient rest.Interface
		ns         string
	}

	func (c *ProjectClient) List(ctx context.Context, opts metav1.ListOptions) (*v3.ProjectList, error) {
		result := v3.ProjectList{}
		err := c.restClient.
			Get().
			Namespace(c.ns).
			Resource("projects").
			VersionedParams(&opts, scheme.ParameterCodec).
			Do(ctx).
			Into(&result)

		return &result, err
	}

	func (c *ProjectClient) Get(ctx context.Context, name string, opts metav1.GetOptions) (*v3.Project, error) {
		result := v3.Project{}
		err := c.restClient.
			Get().
			Namespace(c.ns).
			Resource("projects").
			Name(name).
			VersionedParams(&opts, scheme.ParameterCodec).
			Do(ctx).
			Into(&result)

		return &result, err
	}

	func (c *ProjectClient) Create(ctx context.Context, project *v3.Project) (*v3.Project, error) {
		result := v3.Project{}
		err := c.restClient.
			Post().
			Namespace(c.ns).
			Resource("projects").
			Body(project).
			Do(ctx).
			Into(&result)

		return &result, err
	}

	func (c *ProjectClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
		opts.Watch = true
		return c.restClient.
			Get().
			Namespace(c.ns).
			Resource("projects").
			VersionedParams(&opts, scheme.ParameterCodec).
			Watch(ctx)
	}

	func (c *CattleManagementV3Client) Cluster() ClusterInterface {
		return &ClusterClient{
			restClient: c.restClient,
		}
	}

	type ClusterInterface interface {
		List(ctx context.Context, opts metav1.ListOptions) (*v3.ClusterList, error)
		Get(ctx context.Context, name string, options metav1.GetOptions) (*v3.Cluster, error)
		Create(ctx context.Context, project *v3.Cluster) (*v3.Cluster, error)
		Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	}

	type ClusterClient struct {
		restClient rest.Interface
	}

	func (c *ClusterClient) List(ctx context.Context, opts metav1.ListOptions) (*v3.ClusterList, error) {
		result := v3.ClusterList{}
		err := c.restClient.
			Get().
			Resource("clusters").
			VersionedParams(&opts, scheme.ParameterCodec).
			Do(ctx).
			Into(&result)

		return &result, err
	}

	func (c *ClusterClient) Get(ctx context.Context, name string, opts metav1.GetOptions) (*v3.Cluster, error) {
		result := v3.Cluster{}
		err := c.restClient.
			Get().
			Resource("clusters").
			Name(name).
			VersionedParams(&opts, scheme.ParameterCodec).
			Do(ctx).
			Into(&result)

		return &result, err
	}

	func (c *ClusterClient) Create(ctx context.Context, project *v3.Cluster) (*v3.Cluster, error) {
		result := v3.Cluster{}
		err := c.restClient.
			Post().
			Resource("clusters").
			Body(project).
			Do(ctx).
			Into(&result)

		return &result, err
	}

	func (c *ClusterClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
		opts.Watch = true
		return c.restClient.
			Get().
			Resource("clusters").
			VersionedParams(&opts, scheme.ParameterCodec).
			Watch(ctx)
	}
*/
type Cluster struct {
	Name        string    `json:"name"`
	DisplayName string    `json:"displayName"`
	UUID        string    `json:"uuid"`
	Projects    []Project `json:"projects"`
}

type Project struct {
	Name        string `json:"name"`
	UUID        string `json:"uuid"`
	DisplayName string `json:"displayName"`
}

type KubeClient struct {
	Kubeconfig string
	Debug      bool
	client     *management.Client
	context    context.Context
	age        time.Time
	lastResult []Cluster
}

func (kube *KubeClient) CetClusters() ([]Cluster, error) {
	var err error = nil
	if kube.client == nil {
		if kube.Debug {
			log.Println("@D No client defined, creating new client")
		}
		err = kube.newConfig()
		if err != nil {
			log.Println("@E Errer creating client configuration")
			kube.client = nil
			return nil, err
		}
		kube.lastResult, err = kube.getClusters()
		if err != nil {
			log.Println("@E Errer Getting cluster data, resetting client")
			kube.client = nil
			return nil, err
		}
		kube.age = time.Now()
	}
	if time.Now().Sub(kube.age).Seconds() > 5 {
		kube.lastResult, err = kube.getClusters()
		if err != nil {
			log.Println("@E Errer Getting cluster data, resetting client")
			kube.client = nil
			return nil, err
		}
		kube.age = time.Now()
	}
	return kube.lastResult, err
}
func (kube *KubeClient) newConfig() error {

	if _, err := os.Stat(kube.Kubeconfig); err != nil {
		kube.Kubeconfig = ""
	}
	var config *rest.Config
	var err error
	if kube.Kubeconfig != "" {
		log.Printf("@I Using kubeconfig in: %v\n", kube.Kubeconfig)
		config, err = clientcmd.BuildConfigFromFlags("", kube.Kubeconfig)
		if err != nil {
			return err
		}
	} else {
		config, err = rest.InClusterConfig()
		log.Println("@I Using in Cluster Configuration")
		if err != nil {
			return err
		}
	}
	clientset, err := newClientFor(config)
	if err != nil {
		return err
	}
	v3.AddToScheme(scheme.Scheme)
	kube.context = context.Background()
	kube.client = clientset
	return nil
}

func (kube *KubeClient) getClusters() ([]Cluster, error) {
	if kube.Debug {
		log.Println("@D getClusters: ")
	}
	clusters, err := kube.client.Cluster.List(&types.ListOpts{})
	if err != nil {
		return nil, err
	}
	clusterList := make([]Cluster, 0)
	for _, cluster := range clusters.Data {
		currentCluster := Cluster{Name: cluster.Name, UUID: cluster.UUID, Projects: []Project{}}
		if kube.Debug {
			log.Printf("@D - %+v \n", cluster)
		}
		projects, err := kube.client.Project.List(&types.ListOpts{Filters: map[string]interface{}{"namespace": currentCluster.Name}})
		//kube.client.Project(currentCluster.Name).List(kube.context, metav1.ListOptions{})
		if err != nil {
			return nil, err
		}

		for _, project := range projects.Data {
			currentProject := Project{Name: project.Name, UUID: project.UUID}
			if kube.Debug {
				log.Printf("@D      - %+v\n", project)
			}
			currentCluster.Projects = append(currentCluster.Projects, currentProject)
		}
		clusterList = append(clusterList, currentCluster)
	}
	return clusterList, nil
}
