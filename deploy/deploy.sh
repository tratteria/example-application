export ENABLE_TRATS="true"

cd spire
chmod +x install.sh
./install.sh
cd ..

if [ "$ENABLE_TRATS" = "true" ]; then
    cd tconfigd/installation
    chmod +x install.sh
    ./install.sh
    cd ../..
fi

cd alpha-stocks-dev
chmod +x deploy.sh
./deploy.sh
cd ..
