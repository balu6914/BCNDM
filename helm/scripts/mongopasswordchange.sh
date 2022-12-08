#########################################################################################################################
# Prequisite      :                                                                                       				#
# Access to the kubernetes cluster where password change has to be done. 								  				#
# download the datapace helm chart from gitlab (ex: 	#https://artifactory-espoo1.int.net.nokia.com:443/artifactory/mn-hc-datapace-helm-local/datapace-xx.x-SNAPSHOT_xxx.tgz  #
# kubectl installation in the local/working machine.                                                      				#
# --------------------------------------------------------------------------------------------------------				#
# Usage           : mongopasswordchange.sh -n <namespace> -f <filepath>                                   				#
# namespace       : Namespace where ndm is deloyed. Ex. ndm                                               				#
# filepath        : Yaml filepath where credentials are stored in datapace. Ex./tmp/datapace/values.yaml  				#
# Example Command : mongopasswordchange.sh -n ndm -f /tmp/datapace/values.yaml                            				#
#########################################################################################################################


while getopts "n:f:" opt
do
   case "$opt" in
      n ) namespace="$OPTARG" ;;
      f ) filepath="$OPTARG" ;;
   esac
done

## Downloads yq for yaml file reading
export VERSION=v4.9.6
export BINARY=yq_linux_386
wget https://github.com/mikefarah/yq/releases/download/${VERSION}/${BINARY} -O ./yq && chmod +x ./yq

## Gets the list of all credentials from yaml file for processing
usernamelist=($(cat $filepath | ./yq e '.mongodb.auth.usernames[]' -))
dblist=($(cat $filepath | ./yq e '.mongodb.auth.databases[]' -))
passlist=($(cat $filepath | ./yq e '.mongodb.auth.passwords[]' -))

echo ${#dblist[@]}
lengthoflist=`expr ${#dblist[@]} - 1`
echo $lengthoflist

## JS file creation for parsing to the mongo DB pod for changing passwords
rootpassword=`kubectl get secret --namespace $namespace dpd-datapace-mongodb -o jsonpath="{.data.mongodb-root-password}" | base64 --decode`
echo 'use admin' > mongoaccess.js
echo 'db.auth("root", "'$rootpassword'")' >> mongoaccess.js
for i in `seq 0 $lengthoflist`;
do
   echo 'use '${dblist[$i]} >> mongoaccess.js
   echo 'db.changeUserPassword("'${usernamelist[$i]}'","'${passlist[$i]}'")' >> mongoaccess.js
done

## changing the credentials for all the passed users in the yaml file.
mongopod=`kubectl get po -n $namespace | grep mongo | awk '{print $1}'`
kubectl exec -it -n $namespace $mongopod mongo < mongoaccess.js


## Removes the temporary credentials file created.
rm mongoaccess.js



