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

```java
// 替代 finalize() 的更安全方式
public class ResourceCleaner {
    private static final ReferenceQueue<Object> queue = new ReferenceQueue<>();
    private static final Set<CleanerRef> refs = ConcurrentHashMap.newKeySet();
    
    // 資源清理執行緒
    static {
        Thread cleanerThread = new Thread(() -> {
            while (true) {
                try {
                    CleanerRef ref = (CleanerRef) queue.remove(); // 阻塞等待
                    ref.clean();
                    refs.remove(ref);
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
        refs.add(ref);
    }
    
    private static class CleanerRef extends PhantomReference<Object> {
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

// 使用方式
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