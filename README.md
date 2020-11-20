# NinjaDog-OM
Open-match components of NinjaDog


# QuickStart
**前置需求**
Open-match安裝: [Open-match core Installation](https://openmatch.dev/site/docs/installation/)
Agones安裝: [Agones core Installation](https://agones.dev/site/docs/installation/)


**1.建立所有元件 image**
```
docker build -t [repsitory-url]/nd-director:[version] ./nd-director
docker build -t [repsitory-url]/nd-frontend:[version] ./nd-frontend
docker build -t [repsitory-url]/nd-matchfunction:[version] ./nd-matchfunction
```
**2.上傳所有元件 image**
```
docker push [repsitory-url]/nd-director:[version]
docker push [repsitory-url]/nd-frontend:[version] 
docker push [repsitory-url]/nd-matchfunction:[version] 
```

**3.部署**
```
kubectl create ns nd-mo-components
kubectl apply -f ./nd-om.yaml
```
