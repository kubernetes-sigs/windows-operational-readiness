# rm -rf kubernetes
# rm ./e2e.test
# git clone https://github.com/kubernetes/kubernetes.git
cd kubernetes
# make WHAT="test/e2e/e2e.test"
# mkdir -p ../e2e_test_binary/darwin/
# cp ./_output/bin/e2e.test ../e2e_test_binary/darwin/e2e.test

make WHAT="test/e2e/e2e.test"
mkdir -p ../e2e_test_binary/linux/
cp ./_output/bin/e2e.test ../e2e.test

# ./build/run.sh make WHAT="test/e2e/e2e.test" KUBE_BUILD_PLATFORMS=linux/amd64
# mkdir -p ../e2e_test_binary/linux/
# cp ./_output/dockerized/bin/linux/amd64/e2e.test ../e2e_test_binary/linux/e2e.test
cd ..
# rm -rf kubernetes
