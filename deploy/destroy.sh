cd alpha-stocks
chmod +x destroy.sh
./destroy.sh
cd ..


cd tratteria/installation
chmod +x uninstall.sh
./uninstall.sh --no-spire
cd ../..

cd spire
chmod +x uninstall.sh
./uninstall.sh
cd ..

