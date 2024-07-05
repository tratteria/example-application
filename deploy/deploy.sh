export ENABLE_TRATS="true"

cd spire
chmod +x install.sh
./install.sh
cd ..

if [ "$ENABLE_TRATS" = "true" ]; then
    cd tratteria/installation
    chmod +x install.sh
    ./install.sh --no-spire
    cd ../..
fi

cd alpha-stocks
chmod +x deploy.sh
./deploy.sh
cd ..
