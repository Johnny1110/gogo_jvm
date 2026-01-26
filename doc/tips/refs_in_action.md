# Soft, Weak, Phantom References 的應用場景

<br>

---

<br>

## SoftReference 實現快取

```java
public class ImageCache {
    private Map<String, SoftReference<Image>> cache = new ConcurrentHashMap<>();
    private ReferenceQueue<Image> queue = new ReferenceQueue<>();
    
    public Image getImage(String path) {
        // 先清理已被回收的項目
        cleanUp();
        
        SoftReference<Image> ref = cache.get(path);
        Image img = (ref != null) ? ref.get() : null;
        
        if (img == null) {
            img = loadImage(path);
            cache.put(path, new SoftReference<>(img, queue));
        }
        
        return img;
    }
    
    private void cleanUp() {
        Reference<? extends Image> ref;
        while ((ref = queue.poll()) != null) {
            // 從 cache 中移除對應的 key
            // 注意：這裡需要額外的映射來找到 key
        }
    }
}
```

<br>

---

<br>

## WeakHashMap 實現

<br>

### `WeakHashMap` 的核心概念: 

* Key 使用 WeakReference 包裝
* 當 Key 被 GC 回收後，整個 Entry 會自動被清理

#### 設計目標是: 

當 key 不再被外部強引用時，讓整個 entry（key-value pair）能自動被清理。這特別適合做「快取」或「額外附加資料」的場景。


<br>


### WeakHashMap 的內部結構

```
┌─────────────────────────────────────────┐
│           WeakHashMap                   │
│  ┌─────────────────────────────────┐    │
│  │  Entry[] table (hash buckets)   │    │
│  └─────────────────────────────────┘    │
│                                         │
│  ┌─────────────────────────────────┐    │
│  │  ReferenceQueue<Object> queue   │    │
│  └─────────────────────────────────┘    │
└─────────────────────────────────────────┘
```

每個 Entry 繼承自 WeakReference<K>：

```
Entry extends WeakReference<K> {
    V value;
    int hash;
    Entry<K,V> next;
}
```

**Entry** 本身就是一個指向 key 的 WeakReference。


### WeakHashMap 運作機制：兩階段清理


#### 階段一：GC 處理 weak reference

當某個 key 不再有外部強引用時：

```
外部程式碼：
    MyKey key = new MyKey("foo");
    map.put(key, someValue);
    key = null;  // 現在沒有強引用指向這個 key 了

GC 發生時：
    1. 發現該 key 物件只剩 WeakReference（Entry）指向它
    2. 將該 key 物件標記為可回收
    3. 把對應的 Entry（WeakReference）加入 ReferenceQueue
    4. Entry.get() 從此返回 null (因為 Entry 是弱指標，他指向的物件已經被標記清理了)
```

<br>

#### 階段二：WeakHashMap 主動清理 stale entries

WeakHashMap 在幾乎每個操作（`get`、`put`、`size` 等）前都會呼叫 `expungeStaleEntries()`：

```java
private void expungeStaleEntries() {
    // 從 queue 中取出所有已經被 GC 處理過的 Entry
    for (Object x; (x = queue.poll()) != null; ) {
        Entry<K,V> e = (Entry<K,V>) x;
        int i = indexFor(e.hash, table.length);
        
        // 從 hash table 中移除這個 entry
        // （省略具體的 linked list 移除邏輯）
        
        e.value = null;  // 重要！讓 value 也能被 GC
    }
}
```

<br>


### 為什麼這樣設計？

#### 問題：為什麼不是 value 用 WeakReference？

如果是 value 用 weak reference，那當 value 被回收時，你會得到一個 key 指向 null 的 entry——這通常沒有意義。

WeakHashMap 的語意是：「我想在某個物件（key）還活著的時候，附加一些資料（value）給它。當那個物件死了，附加資料也應該消失。」


#### 問題：為什麼需要 ReferenceQueue？

沒有 ReferenceQueue 的話，你只能透過「遍歷整個 table，檢查每個 entry.get() 是否為 null」來清理。這是 O(n) 操作。

有了 ReferenceQueue，GC 會主動通知你哪些 reference 已經失效，清理變成 O(k)，k 是失效的數量。

#### 問題：為什麼清理是 lazy 的？

WeakHashMap 不使用 background thread 或 finalizer 來清理，而是在每次操作時順便清理。這是刻意的設計：

* 避免同步問題：不需要額外的 lock
* 簡單可預測：行為完全由呼叫者控制
* 效能考量：分攤清理成本到每次操作

<br>

### 範例

```java
public class SimpleWeakHashMap<K, V> {
    private Map<WeakKey<K>, V> map = new HashMap<>();
    private ReferenceQueue<K> queue = new ReferenceQueue<>();
    
    public V put(K key, V value) {
        cleanUp();
        return map.put(new WeakKey<>(key, queue), value);
    }
    
    public V get(K key) {
        cleanUp();
        return map.get(new WeakKey<>(key, null));
    }
    
    private void cleanUp() {
        WeakKey<K> ref;
        while ((ref = (WeakKey<K>) queue.poll()) != null) {
            map.remove(ref);
        }
    }
    
    // Key 包裝類
    private static class WeakKey<K> extends WeakReference<K> {
        private final int hash;
        
        WeakKey(K key, ReferenceQueue<K> queue) {
            super(key, queue);
            this.hash = key.hashCode();
        }
        
        @Override
        public int hashCode() {
            return hash;
        }
        
        @Override
        public boolean equals(Object obj) {
            if (this == obj) return true;
            if (!(obj instanceof WeakKey)) return false;
            Object k1 = this.get();
            Object k2 = ((WeakKey<?>) obj).get();
            return k1 != null && k1.equals(k2);
        }
    }
}
```


<br>

---

<br>

## PhantomReference 追蹤資源釋放

當我們需要追蹤某個物件被 GC 清理掉的事件，並觸發一些連帶事件時，使用 `PhantomReference` 是替代 `finalize()` 的更安全方式

```java
// ResourceCleaner 裡面記錄哪些物件被監視，以及當清理事件發生時該被觸發的 clean 行為
public class ResourceCleaner {
    private static final ReferenceQueue<Object> queue = new ReferenceQueue<>();
    private static final Set<CleanerRef> refs = ConcurrentHashMap.newKeySet();
    
    // 資源清理執行緒
    static {
        Thread cleanerThread = new Thread(() -> {
            while (true) {
                try {
                    CleanerRef ref = (CleanerRef) queue.remove(); // 阻塞等待
                    ref.clean();                                  // 執行自定義的 clean() 方法 
                    refs.remove(ref);                             // 從監管列表中除出
                } catch (InterruptedException e) {
                    break;
                }
            }
        });
        cleanerThread.setDaemon(true);
        cleanerThread.start();
    }
    
    // 註冊需要清理的資源
    public static void register(Object obj, Runnable cleanupAction) {
        CleanerRef ref = new CleanerRef(obj, queue, cleanupAction);
        refs.add(ref); // 加入監管列表
    }
    
    // Wrapper for Watched Target Object
    private static class CleanerRef extends PhantomReference<Object> {
        // 自定義的清理 Event 方法
        private final Runnable cleanupAction;
        
        CleanerRef(Object referent, ReferenceQueue<Object> q, Runnable action) {
            super(referent, q);
            this.cleanupAction = action;
        }
        
        void clean() {
            cleanupAction.run();
        }
    }
}
```
`

#### 使用方式

```java
class NativeResource {
    private long nativePtr;
    
    NativeResource() {
        this.nativePtr = allocateNative();
        
        // 註冊清理器
        long ptr = this.nativePtr; // capture for lambda
        ResourceCleaner.register(this, () -> {
            freeNative(ptr);
            System.out.println("Native resource freed!");
        });
    }
}
```