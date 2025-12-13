/**
 * v0.2.8 測試程式：類型系統增強
 *
 * 測試指令：
 * - instanceof (0xC1)
 * - checkcast (0xC0)
 * - anewarray (0xBD)
 *
 * 編譯指令：javac TestTypeSystem.java
 * 執行指令：gogo_jvm TestTypeSystem.class -debug
 */
public class TestTypeSystem {

    public static void main(String[] args) {
        // ============================================
        // Test 1: instanceof 基本測試
        // ============================================
        testInstanceof();

        // ============================================
        // Test 2: checkcast 基本測試
        // ============================================
        testCheckcast();

        // ============================================
        // Test 3: anewarray 基本測試
        // ============================================
        testAnewarray();

        // ============================================
        // Test 4: 繼承關係測試
        // ============================================
        testInheritance();

        System.out.println(999);  // 全部測試通過
    }

    // Test 1: instanceof 基本測試
    public static void testInstanceof() {
        Animal dog = new Dog();

        // dog instanceof Dog → true
        if (dog instanceof Dog) {
            System.out.println(1);  // 應輸出 1
        }

        // dog instanceof Animal → true (父類)
        if (dog instanceof Animal) {
            System.out.println(2);  // 應輸出 2
        }

        // null instanceof Dog → false
        Animal nullAnimal = null;
        if (nullAnimal instanceof Dog) {
            System.out.println(-1);  // 不應該輸出
        } else {
            System.out.println(3);  // 應輸出 3
        }
    }

    // Test 2: checkcast 基本測試
    public static void testCheckcast() {
        Animal dog = new Dog();

        // 向下轉型：Animal → Dog（成功）
        Dog d = (Dog) dog;
        d.bark();  // 應輸出 10

        // null 轉型（成功，不拋異常）
        Animal nullAnimal = null;
        Dog nullDog = (Dog) nullAnimal;  // 不會拋異常
        System.out.println(4);  // 應輸出 4
    }

    // Test 3: anewarray 基本測試
    public static void testAnewarray() {
        // 建立 Animal 陣列
        Animal[] animals = new Animal[3];

        // 放入元素
        animals[0] = new Dog();
        animals[1] = new Cat();
        animals[2] = new Dog();

        // 讀取並呼叫方法
        animals[0].speak();  // 應輸出 100 (Dog.speak)
        animals[1].speak();  // 應輸出 200 (Cat.speak)

        // 檢查陣列長度
        int len = animals.length;
        System.out.println(len);  // 應輸出 3
    }

    // Test 4: 繼承關係測試
    public static void testInheritance() {
        Dog dog = new Dog();

        // Dog instanceof Animal → true
        if (dog instanceof Animal) {
            System.out.println(5);  // 應輸出 5
        }

        // 向上轉型（自動，不需要 checkcast）
        Animal animal = dog;
        animal.speak();  // 應輸出 100 (Dog.speak，多型)
    }
}

// 父類
class Animal {
    public void speak() {
        System.out.println(0);  // 基類輸出 0
    }
}

// 子類：Dog
class Dog extends Animal {
    public void speak() {
        System.out.println(100);  // Dog 輸出 100
    }

    public void bark() {
        System.out.println(10);  // bark 輸出 10
    }
}

// 子類：Cat
class Cat extends Animal {
    public void speak() {
        System.out.println(200);  // Cat 輸出 200
    }
}