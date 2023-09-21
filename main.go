package main

import (
	"flag"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	v3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

const (
	GroupName    = "management.cattle.io"
	GroupVersion = "v3"
)

var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: GroupVersion}

var (
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme   = SchemeBuilder.AddToScheme
)

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&v3.Project{},
		&v3.ProjectList{},
		&v3.Cluster{},
		&v3.ClusterList{},
	)

	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}

type CattleManagementV3Interface interface {
	Project(namespace string) ProjectInterface
	Projects(namespace string) ProjectInterface
	Cluster() ProjectInterface
	Clusters() ProjectInterface
}

type CattleManagementV3Client struct {
	restClient rest.Interface
}

func NewForConfig(c *rest.Config) (*CattleManagementV3Client, error) {
	config := *c
	config.ContentConfig.GroupVersion = &schema.GroupVersion{Group: GroupName, Version: GroupVersion}
	config.APIPath = "/apis"
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &CattleManagementV3Client{restClient: client}, nil
}

func (c *CattleManagementV3Client) Projects(namespace string) ProjectInterface {
	return &projectClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}

type ProjectInterface interface {
	List(opts metav1.ListOptions) (*v3.ProjectList, error)
	Get(name string, options metav1.GetOptions) (*v3.Project, error)
	Create(*v3.Project) (*v3.Project, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
}

type projectClient struct {
	restClient rest.Interface
	ns         string
}

func (c *projectClient) List(opts metav1.ListOptions) (*v3.ProjectList, error) {
	result := v3.ProjectList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("projects").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

func (c *projectClient) Get(name string, opts metav1.GetOptions) (*v3.Project, error) {
	result := v3.Project{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("projects").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

func (c *projectClient) Create(project *v3.Project) (*v3.Project, error) {
	result := v3.Project{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("projects").
		Body(project).
		Do().
		Into(&result)

	return &result, err
}

func (c *projectClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource("projects").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	AddToScheme(scheme.Scheme)
	crdConfig := *config
	crdConfig.ContentConfig.GroupVersion = &schema.GroupVersion{Group: GroupName, Version: GroupVersion}
	crdConfig.APIPath = "/apis"
	crdConfig.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}
	crdConfig.UserAgent = rest.DefaultKubernetesUserAgent()

	exampleRestClient, err := rest.UnversionedRESTClientFor(&crdConfig)

	result := v3.ProjectList{}
	err := exampleRestClient.
		Get().
		Resource("clusters").
		Do().
		Into(&result)
}
