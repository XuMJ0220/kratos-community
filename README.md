系统采用微服务架构，由一个 API 网关和多个后端业务服务组成。服务间通过 gRPC 进行同步通信，通过 Kafka 进行异步解耦。所有服务均实现容器化部署，并通过 Consul 进行服务治理和配置管理。
graph TD
    subgraph "用户端"
        A[Client App]
    end

    subgraph "基础设施"
        C[Consul]
        K[Kafka]
        P[Prometheus] --> G[Grafana]
    end

    subgraph "数据存储"
        M[MySQL]
        R[Redis]
    end

    subgraph "后端微服务 (Monorepo)"
        GW[API Gateway]

        subgraph "业务服务"
            U[User Service]
            T[Content Service]
            I[Interaction Service]
            L[Relation Service]
        end
    end

    A -- HTTP/REST --> GW
    GW -- gRPC --> U
    GW -- gRPC --> T
    GW -- gRPC --> I
    GW -- gRPC --> L

    U & T & I & L -- "注册/发现/拉取配置" --> C
    U & T & I & L -- "读写数据" --> M & R
    T -- "发布事件" --> K
    U & T & I & L -- "上报监控/健康状态" --> P & C
