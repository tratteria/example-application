cd alpha-stocks-dev
chmod +x destroy.sh
./destroy.sh
cd ..


cd tconfigd/installation
chmod +x uninstall.sh
./uninstall.sh --no-spire
cd ../..

cd spire
chmod +x uninstall.sh
./uninstall.sh
cd ..

