## Usage

[Helm](https://helm.sh) must be installed to use the charts.  Please refer to
Helm's [documentation](https://helm.sh/docs) to get started.

Once Helm has been set up correctly, add the repo as follows:

    helm repo add cert-manager-webhook-etcd http://Blarc.github.io/cert-manager-webhook-etcd

If you had already added this repo earlier, run `helm repo update` to retrieve
the latest versions of the packages.  You can then run `helm search repo
{alias}` to see the charts.

To install the cert-manager-webhook-etcd chart:

    helm install demo cert-manager-webhook-etcd/cert-manager-webhook-etcd

To uninstall the chart:

    helm delete demo
