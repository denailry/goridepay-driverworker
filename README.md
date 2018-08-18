**GET** _/get-order-list?driverID=DRIVER_ID_
```json
[
    {
        "OrderID": 123,
        "Origin": "Jl. Kelud Kiri Atas, Jatibening Baru",
        "Destination": "Jl. Kelud Kanan Atas, Jatibening Baru",
        "DestinationDistance": 2,
        "OriginDistance": 3
    },
        {
        "OrderID": 123,
        "Origin": "Jl. Kelud Kiri Atas, Jatibening Baru",
        "Destination": "Jl. Kelud Kanan Atas, Jatibening Baru",
        "DestinationDistance": 2,
        "OriginDistance": 3
    }
]
```

**POST** _/order_
```json
{
    "OrderID": 1,
    "Origin": "Jl. Kelud Kiri Atas, Jatibening Baru",
    "Destination": "Jl. Kelud Kanan Atas, Jatibening Baru",
    "DestinationDistance": 10,
    "TransactionID": 123,
    "DriverData": 
        [
            {
                "DriverID": 5,
                "OriginDistance": 4
            },
            {
                "DriverID": 2,
                "OriginDistance": 8
            },
            {
                "DriverID": 10,
                "OriginDistance": 3
            },
            {
                "DriverID": 5,
                "OriginDistance": 9
            }
        ]
}
```

**POST** _/invalidate_
```json
{
    "OrderID": 1
}
```