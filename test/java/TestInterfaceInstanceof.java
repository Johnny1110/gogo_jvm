// ============================================================
// TestInterfaceInstanceof.java - 介面與 instanceof 測試
// ============================================================
// 測試目標：
//   1. 物件 instanceof 介面
//   2. 介面引用 instanceof 類
interface Flyable {
    void fly();
}

class Bird implements Flyable {
    public void fly() {
        System.out.println("bird flying");
    }
}

public class TestInterfaceInstanceof {
    public static void main(String[] args) {
        Bird bird = new Bird();

        // bird instanceof Flyable → true
        if (bird instanceof Flyable) {
            System.out.println("bird is a Flyable Object");
        }

        Flyable f = bird;
        // f instanceof Bird → true
        if (f instanceof Bird) {
            System.out.println("f is a bird");
        }

        System.out.println("Done");
    }
}