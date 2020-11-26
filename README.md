# NinjaDog-OM
Open-match components of NinjaDog


# QuickStart
**0.前置需求** <br/>
Open-match安裝: [Open-match core Installation](https://openmatch.dev/site/docs/installation/) <br/>
Agones安裝: [Agones core Installation](https://agones.dev/site/docs/installation/) <br/>


**1.建立所有元件 image**
```
docker build -t [repository-url]/nd-director:[version] ./nd-director
docker build -t [repository-url]/nd-frontend:[version] ./nd-frontend
docker build -t [repository-url]/nd-matchfunction:[version] ./nd-matchfunction
```
**2.上傳所有元件 image**
```
docker push [repository-url]/nd-director:[version]
docker push [repository-url]/nd-frontend:[version] 
docker push [repository-url]/nd-matchfunction:[version] 
```

**3.部署**
```
kubectl create ns nd-mo-components
kubectl apply -f ./nd-om.yaml
```


# Auto-Deployment
修改.github/workflows/google.yml檔案，更改自動部署的目標Cluster
```
env:
  PROJECT_ID: cheep2workshop
  GKE_CLUSTER: [目標Cluster]
  GKE_ZONE: asia-east1-a 	   
  IMAGE_REPO: "asia.gcr.io/cheep2workshop"
```

****
# TroubleShooting
**1.ServiceAccount**<br
當遇到ServiceAccount的權限問題時，使用directRBAC.yaml建立ClusterRole
```
kubectl apply -f directRBAC.yaml
```
綁定ClusterRole至[namespace]:[serviceAccount] (預設為default:default)
```
kubectl create clusterrolebinding [bindingName] --serviceAccount=default:default --clusterrole=nd-director-role
```
