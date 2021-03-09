docker postgres 설치
docker pull postgres:12-alpine


no matching manifest for windows/amd64 10.0.19041 in the manifest list entries
발생 시

Settings - Docker Engine - 
```json
{
  "registry-mirrors": [],
  "insecure-registries": [],
  "debug": false,
  "experimental": false
}
```
experimental를 true로 바꾸고 적용 후 restart

다시 실행한다.
docker pull postgres:12-alpine

docker run --name <container_name> -e <environment_variable> -p <host_ports:container_ports> -d <image>:<tag> 
docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=1234 -d postgres:12-alpine
  
Linux or PowerShell
docker exec -it postgres12 /bin/sh

git bash
winpty docker exec -it postgres12 bash

createdb --username=root --owner=root mytest_db
psql mytest_db

잘 접속되었는지 테스트해 본다.
```bash
mytest_db=# select now();
              now
-------------------------------
 2021-03-07 09:48:31.256488+00
(1 row)

mytest_db=# \q
bash-5.1# dropdb mytest_db
bash-5.1# exit
```

윈도우에서 make 명령어를 사용하기 위해 아래 링크 접속하여 설치
http://gnuwin32.sourceforge.net/packages/make.htm


환경변수 등록
ex) C:\Program Files (x86)\GnuWin32\bin

잘 등록되었는지 확인
```bash
$ make -v
GNU Make 3.81
Copyright (C) 2006  Free Software Foundation, Inc.
This is free software; see the source for copying conditions.
There is NO warranty; not even for MERCHANTABILITY or FITNESS FOR A
PARTICULAR PURPOSE.

This program built for i386-pc-mingw32
```


```bash
$ winpty docker exec -it postgres12 createdb --username=root --owner=root mytest_db
$ winpty docker exec -it postgres12 psql -U root mytest_db
```

https://github.com/golang-migrate/migrate/tree/master/cmd/migrate
이 링크에서 migrate 설치한 후

scoop 사용하여 migrate 설치하기(윈도우 기준)
https://scoop.sh/
scoop 설치
```bash
iwr -useb get.scoop.sh | iex
```

아래와 같이 설치 에러 발생 시
![image](https://user-images.githubusercontent.com/30817924/110441620-5240ed00-80fd-11eb-811d-28ae577448b1.png)
```bash
Set-ExecutionPolicy RemoteSigned -scope CurrentUser
```

![image](https://user-images.githubusercontent.com/30817924/110442157-eb700380-80fd-11eb-9ebc-9c6c457eedd0.png)
```bash
set-executionpolicy -s cu unrestricted
```

아래 그림처럼 설치는 되었다고 나오지만 scoop 명령어가 동작하지 않을 때
![image](https://user-images.githubusercontent.com/30817924/110446396-77842a00-8102-11eb-9a0a-5636f6a8d702.png)

ex) C:\Users\<username> 경로에 scoop 폴더가 있는데 이 폴더를 삭제하고 다시 설치하면 된다.

```bash
cd migrate create -ext sql -dir db-migration -seq init_schema
```
