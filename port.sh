kill $(lsof -i :3000 | grep main | awk '{print $2}' | head -n 1)
echo "air was killed"
