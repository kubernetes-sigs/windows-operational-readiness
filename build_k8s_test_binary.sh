# rm -rf kubernetes
# rm ./e2e.test
# git clone https://github.com/kubernetes/kubernetes.git
cd kubernetes
./build/run.sh make WHAT="test/e2e/e2e.test" KUBE_BUILD_PLATFORMS=linux/amd64
cp ./_output/dockerized/bin/linux/amd64/e2e.test ../e2e.test
cd ..
# rm -rf kubernetes
