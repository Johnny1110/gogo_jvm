# GOGO JVM

<br>

---

<br>


## Process

1. 讀取 .class 文件
2. 使用 ClassLoader 加載類
3. 找到 main 方法
4. 執行解釋器

<br>

目前實現到方法區，使用 ClassLoader 來載入類別，並使用解釋器直接執行 main 方法。

<br>

### 編譯完成後執行:

run class without debug mode:
```bash
./gogo_jvm TargetClass   
```

<br>

run class with debug mode:
```bash
./gogo_jvm TargetClass -debug
```