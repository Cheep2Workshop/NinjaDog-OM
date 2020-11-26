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

****
# Desription
NinjaDog-OM總共有三個元件分別執行不同的功能，分別是nd-frontend、nd-matchfunction、nd-director <br/>

**1.nd-frontend** <br/>
與客戶端的接口
接收客戶端排隊請求，產生Ticket(用於辨識客戶端)，並發送Ticket至Open-match frontend；最後等待回傳的遊戲伺服器資訊，將相對應的Ticket從隊列中刪除。

**2.nd-matchfunction** <br/>
排隊主邏輯
從Open-match query取得目前排隊的玩家Ticket，根據自訂的排隊邏輯(依Rank, Mode, etc...)將玩家群分成多個集合(Match)，並回傳給Open-match backend。

**3.nd-director** <br/>
Open-match與Agones對接的元件
從Open-match backend取得來自nd-matchfunction的結果(Match)，向Agones-Allocator申請分配遊戲伺服器(Allocation)，分析遊戲伺服器資訊(Assignment含IP和Port)並回傳給Open-match backend。

****
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
**1.ServiceAccount** <br/>
當遇到ServiceAccount的權限問題時，例
```
User "system:serviceaccount:default:default" cannot create services in the namespace "default"
```
首先使用directRBAC.yaml建立ClusterRole
```
kubectl apply -f directRBAC.yaml
```
其次產生ClusterRoleBinding並綁定ClusterRole至ServiceAccount (預設為default:default)
```
kubectl create clusterrolebinding [bindingName] --serviceAccount=default:default --clusterrole=nd-director-role
```
