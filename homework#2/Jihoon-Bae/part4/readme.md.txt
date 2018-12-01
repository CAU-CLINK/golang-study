비트코인 트랜잭션.

트랜잭션은 입력과 출력의 조합이다.

![](https://mingrammer.com/images/2018-05-26-transactions-diagram.png)
트랜잭션을 이해하기위한 아주 좋은 그림.

트랜잭션은 스크립트를 사용해 값을 잠그기만 하며, 스크립트로 잠근 사람만 잠금을 해제할 수 있다.

-. 비트코인은 내부적으로 출력 잠금과 로직 해제를 정의하는데 사용되는 *Script*라는 스크립팅 언어를 사용한다. 이 언어는 아주 원시적이지만 (해킹과 잘못된 사용을 피하기위해 의도된 설계이다), 여기서 자세히 논의하지는 않는다

-ScriptSig는 출력의 ScriptPubKey에서 사용되는 데이터를 제공하는 스크립트이다. 데이터가 올바르면 출력의 잠금을 해제할 수 있으며 출력의 값 (Value)을 사용해 새로운 출력을 생성할 수 있다. 그렇지 않은 경우에 해당 출력은 입력에서 참조할 수 없다.

트랜잭션 추가했으니 바뀌어야할 것들
-newcoinbase()
-block struct에 tx추가
-newblock(), newgenesisBlock()에서 data->tx로수정
-createBlockchain(): cbtx(코인베이스tx)로 제네시스블락생성,newblockchain()
-pow prepare()수정, HashTransaction() : 모든 tx들합해서 단일hash화.

미사용 트랜잭션들을 출력해보자
-UTXO를 찾아야한다.
-키로 잠금해제 가능한 tx들만 찾으면 됨
