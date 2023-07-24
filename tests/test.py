import requests

BASE_URL = "http://localhost:3333"


def test_e2e():

    # Register two donations
    result = requests.post(
        BASE_URL + "/donations/register",
        json={
            "name": "Andrew",
            "type": "money",
            "quantity": 55,
            "description": "hello",
        },
    )
    result_json = result.json()
    assert result_json["name"] == "Andrew"
    assert result_json["type"] == "money"
    assert result_json["quantity"] == 55
    assert result_json["description"] == "hello"
    assert result_json["id"] == 1

    result = requests.post(
        BASE_URL + "/donations/register",
        json={
            "name": "Jerry",
            "type": "clothing",
            "quantity": 11,
            "description": "socks",
        },
    )
    result_json = result.json()
    assert result_json["name"] == "Jerry"
    assert result_json["type"] == "clothing"
    assert result_json["quantity"] == 11
    assert result_json["description"] == "socks"
    assert result_json["id"] == 2

    # Make a few distributions
    result = requests.post(
        BASE_URL + "/donations/distribute",
        json={"donation_id": 1, "type": "money", "quantity": 54, "description": "yes"},
    )
    result_json = result.json()
    assert result_json["donation_id"] == 1
    assert result_json["type"] == "money"
    assert result_json["quantity"] == 54
    assert result_json["description"] == "yes"
    assert result_json["id"] == 1

    result = requests.post(
        BASE_URL + "/donations/distribute",
        json={"donation_id": 1, "type": "money", "quantity": 1, "description": "yes"},
    )
    result_json = result.json()
    assert result_json["donation_id"] == 1
    assert result_json["type"] == "money"
    assert result_json["quantity"] == 1
    assert result_json["description"] == "yes"
    assert result_json["id"] == 2

    result = requests.post(
        BASE_URL + "/donations/distribute",
        json={
            "donation_id": 2,
            "type": "clothing",
            "quantity": 10,
            "description": "socks",
        },
    )
    result_json = result.json()
    assert result_json["donation_id"] == 2
    assert result_json["type"] == "clothing"
    assert result_json["quantity"] == 10
    assert result_json["description"] == "socks"
    assert result_json["id"] == 3

    # Make a bad distribution
    # Go over quantity
    result = requests.post(
        BASE_URL + "/donations/distribute",
        json={"donation_id": 1, "type": "money", "quantity": 1, "description": "yes"},
    )
    assert result.status_code != 200

    # Type doesn't match donation id
    result = requests.post(
        BASE_URL + "/donations/distribute",
        json={"donation_id": 2, "type": "money", "quantity": 1, "description": "socks"},
    )
    assert result.status_code != 200

    result = requests.get(BASE_URL + "/donations/report/inventory")
    result_json = result.json()
    assert len(result_json) == 3
    assert len(result_json["money"]) == 1
    assert len(result_json["clothing"]) == 1
    assert result_json["money"][0]["donation"]["id"] == 1
    assert len(result_json["money"][0]["distributions"]) == 2
    assert len(result_json["clothing"]) == 1
    assert result_json["clothing"][0]["donation"]["id"] == 2
    assert len(result_json["clothing"][0]["distributions"]) == 1

    result = requests.get(BASE_URL + "/donations/report/donors")
    result_json = result.json()
    print(result_json)
    assert len(result_json["report"]) == 2
    for d in result_json["report"]:
        if d["name"] == "Andrew":
            assert d["donations"]["money"]["quantity"] == 55
            assert d["donations"]["money"]["quantity_distributed"] == 55
        else:
            assert d["donations"]["clothing"]["quantity"] == 11
            assert d["donations"]["clothing"]["quantity_distributed"] == 10
