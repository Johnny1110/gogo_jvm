# Gogo JVM

<br>

---

<br>


Using GoLang to implement a simple JVM

<br>

## MVP plan

* [MVP plan](doc/mvp_plan)

<br>
<br>

## Structure

<br>

1. [JVM Entry](cmd/gogo_jvm)
2. [Classfile](classfile)
3. [JVM Runtime Data Area](runtime)
   * [Method-Area](runtime/method_area)
   * [Heap](runtime/heap)
   * [Slot](runtime/rtcore)
4. [instructions](instructions)
5. [interpreter](interpreter)
6. [native](native)



### Tips:

* Java Constant Pools: [link](doc/tips/constant_pool.md)
* Java Method Descriptor: [link](doc/tips/descriptor.md)
* ClassLoader: 為什麼類別載入時的 static field slot ID 計算時不需要考慮父類別？: [link](doc/tips/class_loader_cal_field_slot_id.md)

