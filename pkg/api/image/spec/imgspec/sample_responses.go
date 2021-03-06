package imgspec

var ImageHistoryGood = `
{
  "data": [
    {
      "Id": "83364c85cafc",
      "CreatedAt": "(.+)",
      "CreatedBy": "",
      "RepoTags": [
        "andreaskokkalis/dc:0.0_seed"
      ],
      "Size": 0,
      "Comment": ""
    }
  ]
}
`
var ImageListGood = `
{
  "data": [
    {
      "Id": "([A-Fa-f0-9]{12,64})$",
      "RepoTags": [
        "andreaskokkalis/dc:testxxx"
      ],
      "CreatedAt": "(.+)"
    },
    {
      "Id": "([A-Fa-f0-9]{12,64})$",
      "RepoTags": [
        "andreaskokkalis/dc:0.3_tcpdump_airplane"
      ],
      "CreatedAt": "(.+)"
    },
    {
      "Id": "([A-Fa-f0-9]{12,64})$",
      "RepoTags": [
        "andreaskokkalis/dc:0.2_tcpdump_assignment"
      ],
      "CreatedAt": "(.+)"
    },
    {
      "Id": "([A-Fa-f0-9]{12,64})$",
      "RepoTags": [
        "andreaskokkalis/dc:0.1_traceroute"
      ],
      "CreatedAt": "(.+)"
    },
    {
      "Id": "([A-Fa-f0-9]{12,64})$",
      "RepoTags": [
        "andreaskokkalis/dc:0.0_seed"
      ],
      "CreatedAt": "(.+)"
    }
  ]
}
`
