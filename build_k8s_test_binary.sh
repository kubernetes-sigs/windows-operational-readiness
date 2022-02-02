git clone https://github.com/kubernetes/kubernetes.git
cd kubernetes
make WHAT="test/e2e/e2e.test"
cp ./_output/bin/e2e.test ../e2e.test
cd ..
rm -rf kubernetes
