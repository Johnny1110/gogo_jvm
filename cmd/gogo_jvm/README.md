# GOGO JVM v0.2.3

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