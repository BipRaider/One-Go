#!/bin/bash
# Change to the directory with our code that we plan to work from
cd "$GOPATH/src/bipgo.pw"
echo "==== Releasing bipgo.pw ===="
echo " Deleting the local binary if it exists (so it isn't uploaded)..."
rm bipgo.pw
echo " Done!"
echo " Deleting existing code..."
ssh root@bipus.bipgo.pw "rm -rf /root/go/src/bipgo.pw/"
echo " Code deleted successfully!"
echo " Uploading code..."
# The \ at the end of the line tells bash that our
# command isn't done and wraps to the next line.
rsync -avr --exclude '.git/*' --exclude 'tmp/*' --exclude 'images/*' ./ root@bipus.bipgo.pw:/root/go/src/bipgo.pw/
#  почемуто папка  /root/go/src/bipgo.pw/?   создаёться
echo " Code uploaded successfully!"
echo " Go getting deps..."
ssh root@bipus.bipgo.pw "go get golang.org/x/crypto/bcrypt"
ssh root@bipus.bipgo.pw "go get github.com/gorilla/mux"
ssh root@bipus.bipgo.pw "go get github.com/gorilla/schema"
ssh root@bipus.bipgo.pw "go get -u github.com/go-sql-driver/mysql"
ssh root@bipus.bipgo.pw "go get github.com/jinzhu/gorm"
ssh root@bipus.bipgo.pw "go get github.com/gorilla/csrf"
echo " Building the code on remote server..."
ssh root@bipus.bipgo.pw 'cd /root/app; go build -o ./server /root/go/src/bipgo.pw/*.go'
echo " Code built successfully!"
echo " Moving assets..."
ssh root@bipus.bipgo.pw "cd /root/app; cp -R /root/go/src/bipgo.pw/assets ."
echo " Assets moved successfully!"
echo " Moving views..."
ssh root@bipus.bipgo.pw "cd /root/app; cp -R /root/go/src/bipgo.pw/views ."
echo " Views moved successfully!"
echo " Moving Caddyfile..."
ssh root@bipus.bipgo.pw "cd /root/app; cp /root/go/src/bipgo.pw/Caddyfile ."
echo " Views moved successfully!"
echo " Restarting the server..."
ssh root@bipus.bipgo.pw "sudo service bipgo.pw restart"
echo " Server restarted successfully!"
echo " Restarting Caddy server..."
ssh root@bipus.bipgo.pw "sudo service caddy restart"
echo " Caddy restarted successfully!"
echo "==== Done releasing bipgo.pw ===="