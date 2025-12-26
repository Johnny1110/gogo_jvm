// ============================================================
// TestInterfaceInheritance.java - 介面繼承測試
// ============================================================
// 測試目標：
//   1. 介面繼承介面
//   2. 子介面包含父介面的方法
interface Animal {
    void makeSound();
}

interface Pet extends Animal {
    void play();
}

class Dog implements Pet {
    public void makeSound() {
        System.out.println("Bark!");
    }

    public void play() {
        System.out.println("Playing Ball...");
    }
}

public class TestInterfaceInheritance {
    public static void main(String[] args) {
        Pet pet = new Dog();
        pet.makeSound();
        pet.play();

        Animal animal = pet;
        animal.makeSound();

        System.out.println("Done");
    }
}