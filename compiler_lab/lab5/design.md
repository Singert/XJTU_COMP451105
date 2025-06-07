# 文法规则

```
S' -> Program

Program -> StmtList | StmtList EOF

Stmt -> Func | Decl | id = Expr ; | return Expr ; | Block 
      | if ( Cond ) Stmt | if ( Cond ) Stmt else Stmt 
      | while ( Cond ) Stmt | id MultiIndex = Expr ; 
      | id ( Args ) ; | for ( ForInit Cond ; Expr ) Stmt

Func -> Type id ( Args ) Block

Decl -> Type id = Expr ; | Type id MultiIndex ; 
      | Type id MultiIndex = Expr ; | Type id MultiIndex = InitList ;

Type -> type_kw

Block -> { StmtList }

StmtList -> ε | Stmt | StmtList Stmt

Expr -> Expr + Term | Expr - Term | Term | InitList | id = Expr
Term -> Term * CastExpr | Term / CastExpr | CastExpr

CastExpr -> CastPrefix Factor | Factor
CastPrefix -> ( Type )

Factor -> id ( Args ) | num | float | char | string | id 
        | ( Expr ) | id MultiIndex | - Factor

Args -> NonEmptyArgs | ε
NonEmptyArgs -> Expr | NonEmptyArgs , Expr
              | Type id | Type id = Expr
              | Type id MultiIndex | Type id MultiIndex = Expr
              | NonEmptyArgs , Type id | NonEmptyArgs , Type id = Expr
              | NonEmptyArgs , Type id MultiIndex | NonEmptyArgs , Type id MultiIndex = Expr

IndexList -> ε | Expr | IndexList , Expr

Cond -> Cond && Cond | Cond || Cond | ! Cond 
      | Expr < Expr | Expr > Expr | Expr <= Expr | Expr >= Expr 
      | Expr != Expr | Expr == Expr | ( Cond ) | Expr

MultiIndex -> ε | [ IndexList ] MultiIndex

InitList -> { } | { NonEmptyInitList }
NonEmptyInitList -> Expr | NonEmptyInitList , Expr

ForInit -> Decl | Expr | ε
```