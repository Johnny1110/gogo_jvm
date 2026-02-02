/**
 * v0.4.0 Test - Basic Thread Test
 *
 * 測試基本的執行緒功能：
 * - Thread 創建
 * - Thread.start()
 * - Thread.join()
 * - Thread.sleep()
 *
 * 預期輸出:
 * 1
 * 2
 * 3
 * 99
 */
public class TestThreadBasic {
    public static void main(String[] args) throws InterruptedException {
        System.out.println(1);  // main 執行緒

        // 創建一個新執行緒
        Thread t = new Thread(new Runnable() {
            public void run() {
                try {
                    Thread.sleep(3000L);
                    System.out.println(2);  // 子執行緒
                } catch (InterruptedException ex) {
                    System.out.println("warning !!! InterruptedException in sub thread!");  // 子執行緒
                }

            }
        });

        // 啟動執行緒
        t.start();

        // 等待子執行緒結束
        Thread.sleep(5000L);

        System.out.println(3);  // 回到 main

        System.out.println(99); // 測試完成
    }
}