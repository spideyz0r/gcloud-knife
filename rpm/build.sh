#!/bin/bash
echo "Recreating rpmbuild directory"
rm -rvf /root/rpmbuild/
rpmdev-setuptree
echo "Copying over sources"
cp -rpv /project/gcloud-knife/{go.mod,go.sum,knife.go,main.go} /root/rpmbuild/SOURCES
echo "Building SRPM"
rpmbuild --undefine=_disable_source_fetch -bs /project/gcloud-knife/rpm/gcloud-knife.spec
mkdir -p ~/.config
mv /project/gcloud-knife/copr ~/.config/copr
copr-cli build gcloud-knife /root/rpmbuild/SRPMS/*.src.rpm