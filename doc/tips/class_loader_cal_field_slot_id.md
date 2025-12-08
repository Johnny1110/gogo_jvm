# 為什麼類別的 static field slot ID 計算時不需要考慮婦類別？

<br>

---

<br>
<br>

### 概述

在 class_loader.go 中，loadNonArrayClass 時會對定義好的 class 進行 link 的動作，執行 `prepare()` 邏輯：

```go
// prepare Preparation
// allocate space for static const
func prepare(class *Class) {
	calcInstanceFieldSlotIds(class) // 1
	calcStaticFieldSlotIds(class)   // 2
	allocAndInitStaticVars(class)   // 3
}
```

1. 計算 Non-Static Fields 需要佔用的 Slot 數量大小
2. 計算 Static Fields 需要佔用的 Slot 數量大小
3. 分配靜態空間 (根據 step 2 計算值進行分配)

<br>

目前看上去蠻合理的，只有一個細節要多提一下，為什麼在 loadNonArrayClass 時只分配靜態空間？

原因就是非靜態空間不屬於 class 管轄，那個要留在 new Object 時才需要分配，現階段只需要先算好佔用多大空間方便 new Object 使用就好。

<br>
<br>

### 問題

接著來看重點：靜態 & 非靜態 Fields Slots 空間計算：

```go
// calcInstanceFieldSlotIds calculate instance fields slot ID
func calcInstanceFieldSlotIds(class *Class) {
	slotId := uint(0)
	// from parent
	if class.superClass != nil {
		slotId = class.superClass.instanceSlotCount
	}
	for _, field := range class.fields {
		if !field.IsStatic() {
			field.slotId = slotId
			slotId++
			if field.isLongOrDouble() {
				slotId++
			}
		}
	}

	class.instanceSlotCount = slotId
}

// calcStaticFieldSlotIds calculate static fields slot ID
func calcStaticFieldSlotIds(class *Class) {
	slotId := uint(0)
	// why not include parent class ?
	for _, field := range class.fields {
		if field.IsStatic() {
			field.slotId = slotId
			slotId++
			if field.isLongOrDouble() {
				slotId++
			}
		}
	}
	class.staticSlotCount = slotId
}
```

<br>

可以看到 `calcInstanceFieldSlotIds()` 時會將父類別的非靜態 Fields 佔用數量也算到當前 class 非靜態變數表的容量內。但是 `calcStaticFieldSlotIds()` 則沒有。

難道說 Class 的 Static Fields 是沒有繼承關係的嗎？

事實上，在寫 Java 中能夠從子類調用父類的靜態變數是 java 語法糖。 Instance Fields 跟 Static Fields 的儲存位置完全不同：

<br>

```
┌─────────────────────────────────────────────────────────────────┐
│  Instance Fields vs Static Fields Storage                       │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  Instance Fields: Store in Object                               │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │  Object (Dog Instant)                                   │    │
│  │  ┌─────────────────────────────────────────────────┐    │    │
│  │  │ fields: [age, name]                             │    │    │
│  │  │          ↑     ↑                                │    │    │
│  │  │       from    from                              │    │    │
│  │  │       Animal   Dog                              │    │    │
│  │  └─────────────────────────────────────────────────┘    │    │
│  └─────────────────────────────────────────────────────────┘    │
│                                                                 │
│  Static Fields: Store in Class                                  │
│  ┌──────────────────┐    ┌──────────────────┐                   │
│  │  Animal (Class)  │    │  Dog (Class)     │                   │
│  │  staticVars: [x] │    │  staticVars: [y] │                   │
│  └──────────────────┘    └──────────────────┘                   │
│         ↑                        ↑                              │
│    Animal.x                   Dog.y                             │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

<br>
<br>

### 用具體例子說明

```java
class Animal {
    int age;              // instance field
    static int count;     // static field
}

class Dog extends Animal {
    String name;          // instance field  
    static int dogCount;  // static field
}
```

<br>

### Instance Fields：需要合併到同一個 Object

當你建立 `new Dog()` 時，這個物件必須同時容納 Animal 和 Dog 的所有實例欄位：

```
┌─────────────────────────────────────────────────────────────────┐
│  Dog object = new Dog();                                        │
│                                                                 │
│  object.fields:                                                 │
│  ┌────────┬────────┐                                            │
│  │ slot 0 │ slot 1 │                                            │
│  │  age   │  name  │                                            │
│  │ (from  │ (from  │                                            │
│  │Animal) │  Dog)  │                                            │
│  └────────┴────────┘                                            │
│                                                                 │
│  Dog's instanceSlotCount = 2                                    │
│  - Animal provide 1 slot (age)                                  │
│  - Dog provide 1 slot (name)                                    │
└─────────────────────────────────────────────────────────────────┘
```

* 這就是為什麼要: `slotId = class.superClass.instanceSlotCount`

<br>
<br>

### Static Fields：各 Class 獨立存放

靜態欄位不是存在物件裡，而是存在各自的 Class 結構中：
```
┌─────────────────────────────────────────────────────────────────┐
│  Animal class:                    Dog class:                    │
│  ┌─────────────────┐              ┌─────────────────┐           │
│  │ staticVars:     │              │ staticVars:     │           │
│  │ ┌────────┐      │              │ ┌────────┐      │           │
│  │ │ slot 0 │      │              │ │ slot 0 │      │           │
│  │ │ count  │      │              │ │dogCount│      │           │
│  │ └────────┘      │              │ └────────┘      │           │
│  └─────────────────┘              └─────────────────┘           │
│                                                                 │
│  Animal.staticSlotCount = 1                                     │
│  Dog.staticSlotCount = 1                                        │
└─────────────────────────────────────────────────────────────────┘
```

<br>
<br>

### 「繼承」靜態欄位的真相

**Java 中靜態欄位「可以被繼承」，但這只是語法糖：**

```java
class Animal {
static int count = 10;
}

class Dog extends Animal {
}

// 這兩種寫法效果相同：
System.out.println(Dog.count);    // 輸出 10
System.out.println(Animal.count); // 輸出 10
```

但在 JVM 層面，`Dog.count` 會被編譯器處理成 `Animal.count`：
```
┌─────────────────────────────────────────────────────────────────┐
│  Java 程式碼          Bytecode                                   │
│                                                                 │
│  Dog.count      →    getstatic Animal.count                     │
│                      (編譯器知道 count 定義在 Animal)              │
│                                                                 │
│  Animal.count   →    getstatic Animal.count                     │
└─────────────────────────────────────────────────────────────────┘
```

* 兩者編譯出來的 bytecode 完全一樣

<br>
<br>

### 驗證方式

可以用 `javap -c Dog.class` 查看 bytecode，會發現存取 `Dog.count` 實際上是：
```
getstatic Animal.count : I