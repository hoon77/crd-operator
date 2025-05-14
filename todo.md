# Todo

|Date|do|desc|
|:------:|---|---|
|2025.4.25|kubebuilder init|kubebuilder 프로젝트 셋팅|
|2025.4.29|webapp_type 설정, 코드 정리|Webhook,Metric 등등 정리|
|2025.4.30|deployment, service 생성|crd `webapp`을 통한 리소스 생성|
|2025.5.2|deployment availableReplicas crd Status에 반영|webapp_controller Owns() 설정, deployment 상태가 바뀌면 자동으로 webapp reconcile 호출|
|2025.5.13|configmap 생성|crd `webapp`을 통한 리소스 생성, reflect.DeepEqual|
|2025.5.14|webapp configData 변경 → configmap 업데이트 → deployment rolling update 발생 → pod 재시작|configData hash 값을 deployment annotations에 반영|
