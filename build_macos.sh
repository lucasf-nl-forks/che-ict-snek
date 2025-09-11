GOOS=darwin GOARCH=amd64 go build
find .

# keychain

KEYCHAIN_PATH=$RUNNER_TEMP/build.keychain-db

security create-keychain -p "$MACOS_P12_PASSWORD" "$KEYCHAIN_PATH"
security unlock-keychain -p "$MACOS_P12_PASSWORD" "$KEYCHAIN_PATH"

security import $RUNNER_TEMP/cert.p12 \
  -P "$MACOS_P12_PASSWORD" \
  -A -t cert -f pkcs12 \
  -k "$KEYCHAIN_PATH"

security import $RUNNER_TEMP/cert_app.p12 \
  -P "$MACOS_P12_PASSWORD" \
  -A -t cert -f pkcs12 \
  -k "$KEYCHAIN_PATH"

security list-keychains -s "$KEYCHAIN_PATH"
security set-key-partition-list -S apple-tool:,apple: -s -k "$MACOS_P12_PASSWORD" "$KEYCHAIN_PATH"

# pkg
mkdir -p pkgroot/usr/local/bin
cp snek pkgroot/usr/local/bin

# sign binary
codesign --force --options runtime --sign "$MACOS_P12_APP_NAME" pkgroot/usr/local/bin/snek
xcrun notarytool submit pkgroot/usr/local/bin/snek \
  --apple-id "$MACOS_APPLE_ID" \
  --team-id "$MACOS_TEAM_ID" \
  --password "$MACOS_APPLE_PASSWORD" \
  --wait
xcrun stapler staple pkgroot/usr/local/bin/snek

# build pkg
pkgbuild \
  --root pkgroot \
  --identifier nl.che-ict.snek \
  --version 0.0 \
  --install-location / \
  snek.pkg

productsign --sign "$MACOS_P12_NAME" snek.pkg snek-signed.pkg
xcrun notarytool submit snek-signed.pkg --apple-id "$MACOS_APPLE_ID" --team-id "$MACOS_TEAM_ID" --password "$MACOS_APPLE_PASSWORD" --wait
xcrun stapler staple snek-signed.pkg

# packaging
mkdir dmgroot
cp snek-signed.pkg dmgroot/

hdiutil create -volname "Snek Installer" \
  -srcfolder dmgroot \
  -ov -format UDZO snek.dmg

# signing round 3
codesign --force --sign "$MACOS_P12_APP_NAME" snek.dmg
xcrun notarytool submit snek.dmg \
  --apple-id "$MACOS_APPLE_ID" \
  --team-id "$MACOS_TEAM_ID" \
  --password "$MACOS_APPLE_PASSWORD" \
  --wait
xcrun stapler staple snek.dmg
