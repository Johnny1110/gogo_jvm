// ============================================================
// TestMultiInterface.java - 多介面實現測試
// ============================================================
// 測試目標：
//   1. 一個類實現多個介面
//   2. 透過不同介面引用呼叫方法
interface Runner {
    void run();
}

interface Swimmer {
    void swim();
}

class Athlete implements Runner, Swimmer {
    public void run() {
        System.out.println("Running");
    }

    public void swim() {
        System.out.println("Swimming");
    }
}

public class TestMultiInterface {
    public static void main(String[] args) {
        Athlete a = new Athlete();

        Runner r = a;
        r.run();

        Swimmer s = a;
        s.swim();

        System.out.println("Done");
    }
}