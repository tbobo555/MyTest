為了方便測試, 就都使用GET Method，

有用transaction，其Isolation為Serializable，

因此不必擔心併發會導致投票數據錯誤。

#API

投 yes

https://my-test-269314.appspot.com/vote?op=yes

投 no

https://my-test-269314.appspot.com/vote?op=no

查看投票結果

https://my-test-269314.appspot.com/vote/info

