


### 1. IAM 설정하기
![image](https://user-images.githubusercontent.com/30817924/148495376-939bd340-380f-481d-b7a1-666ff02c3217.png)
- IAM 클릭
   
#### 1.1. 사용자 그룹 생성하기
![image](https://user-images.githubusercontent.com/30817924/148495432-5150ce0f-4143-4bf5-a9fb-b37aade2301e.png)
- 사용자 그룹 클릭 후 그룹 생성 클릭
   
#### 1.2. 사용자 그룹 정보 입력하기
![image](https://user-images.githubusercontent.com/30817924/148495511-1ed944bb-7d08-4ef6-a9d3-14be4bf414df.png)
1. 사용자 그룹 이름 입력
2. 권한 정책 연결 - AmazonEC2ContainerRegistryFullAccess 체크 후 그룹생성
   
#### 1.3. 사용자 그룹 생성 완료
![image](https://user-images.githubusercontent.com/30817924/148495613-e80b3219-ad14-4990-aaea-70a0c6474ded.png)
   
#### 2.1. 사용자 생성하기
![image](https://user-images.githubusercontent.com/30817924/148495648-6103a323-fd73-4ed5-9f84-bda0c5e24fb6.png)
- 왼쪽 사용자 메뉴 클릭 후 사용자 추가 클릭
   
#### 2.2. 사용자 정보 입력하기
![image](https://user-images.githubusercontent.com/30817924/148495678-d47bd06d-8c7b-49bc-b5c8-b6f96e372d31.png)
- 사용자 이름을 입력하고 액세스 키 체크

#### 2.3. 생성할 사용자 그룹 설정하기
![image](https://user-images.githubusercontent.com/30817924/148495742-fab5a7dd-b7c5-4234-9120-4d9271baabbf.png)
- 위에서 만들었던 그룹을 체크하고 다음 클릭

![image](https://user-images.githubusercontent.com/30817924/148495796-24072e89-aefc-4d32-8277-a27af58585dd.png)
- 다음

![image](https://user-images.githubusercontent.com/30817924/148495807-bc87cb6b-9576-4113-8748-2438fd76312b.png)
- 다음

#### 2.4. 액세스 키 확인
![image](https://user-images.githubusercontent.com/30817924/148495818-6c13cc5c-147b-4de9-b450-32484d66a552.png)
- 사용자가 추가되었고 액세스 키 ID와 비밀 액세스 키가 생성되었다. 이 두 가지를 사용하므로 csv 다운로드하여 저장하여 보관한다.(저장하지 않으면 다시 확인할 수 없다.)

### 2. ECR 리포지토리 생성하기

#### 1.1. 리포지토리 생성하기
![image](https://user-images.githubusercontent.com/30817924/148495979-e244e949-8924-4c0d-b71d-ade315b7727c.png)
- ecs 검색하여 접속한다.

![image](https://user-images.githubusercontent.com/30817924/148496001-a511261b-d03c-468d-bda6-4fa3eae165a5.png)
- 왼쪽 메뉴 Amazon ECR 리포지토리 클릭

![image](https://user-images.githubusercontent.com/30817924/148496035-ffc38b7b-302b-4ad3-ba1a-043748e28c02.png)
- 리포지토리 생성 클릭

![image](https://user-images.githubusercontent.com/30817924/148496060-4ecc313b-145f-4032-9d4d-8cd4ef78b4f9.png)
- 프라이빗, 퍼블릭 중 선택하고 리포지토리 이름입력 후 생성한다.
- 리포지토리 이름은 *$ECR_REPOSITORY*에서도 사용한다.

![image](https://user-images.githubusercontent.com/30817924/148498798-e033ce70-0a44-499d-843d-26b55bc4ff9b.png)
- 소스코드 푸시하면 위와 같이 이미지가 생성된다.
