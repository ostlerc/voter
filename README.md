voter
=====

voting strategy analysis software

egen
====
  randomly generate an election

    stdin = election [csv]

    flags
        -cand=3: Number of candidates in election
        -cond=false: Force condorcet winner
        -fix=false: Use fixed random seed
        -o="json": output type [json,csv]
        -peak=false: Force generation with peak preference
        -pref=false: Force preference of some candidate
        -rand=false: Use a random vote/cand count
        -vote=4: Number of voters in election
        -weight=5: Maximum weight of a vote

    stdout: election json
    {
        "votes":
            [ {
                "peak": 0,
                "vote": {
                    "0": 0
                    "1": 2,
                    "2": 1,
                }
            },
            {
                ...
            }, ],

            "peak": true,
            "pref": {
                "b": 1,
                "a": 2
            }
        "rank": [ 4, 3, 5, 0, 0, 0 ],
            "condorcet": 2,
            "candidates": 3
    }

tally
=====
  tallies election results

    stdin = election [csv,json]

    flags
        -i="json": tally input type. [json,csv]
        -o="json": output type [json,csv]
        -t="slater,kemeny": tally type results

    - sample json output
    {
        "names": {
            "0": "Alex"
            "1": "Bart",
            ...
        },
        "results": {
            "slater": [ 0, 3, 4, 1, 5, 2, 6 ],
            "kemeny": [ 4, 0, 3, 5, 2, 1, 6 ],
            ...
        }
    }
    - sample csv output
        rank,slater,kemeny
        1,Alex,Erik
        2,David,Alex
        ...

graph
=====
  majority graph of election

    stdin = election [csv,json]

    flags
    -i="json": graph input type. [csv,json]
    -o="dot": graph output type. [json,dot]

    - sample dot output
    digraph G {
        Alex -> Greg [label="12"];
        Alex -> Bart [label="6"];
        Alex -> Cindy [label="6"];
        ...
    }

    - sample json output
    {
        "nodes": {
            "6": {},
            "5": {
                "edges": {
                    "6": 18,
                    "2": 6
                }
            },
            "4": {
                "edges": {
                    "6": 18,
                    ...
                }
            },
            ....
        }
    }

Examples
========

    cat sample.csv | graph | dot -Tpng | feh -

In this example we pipe an election csv file to the graph binary and visualize the majority graph as a png

    egen -cand 5 -vote 20 -pref

In this example we generate a json election with 5 candidates (seats) and 20 votes forcing a preference.
Note that each 'vote' is actually a weighted vote with default weights ranging from 1-5. See help on
egen for more useful flags

    egen | graph | dot -Tpng | feh -

In this example we generate a new election with all the default flags and create a dot file majority graph

    cat sample.csv | tally -i csv -o csv | column -s, -tn

In this example we tally an election from a csv file and output it in csv format. Then pretty print it.

    egen | tally | jq .

In this example we generate a random election and tally it into a json result. The result is then pretty printed

    cat sample.csv | egen -cond -o csv | column -s, -tn

This example generates an election from a csv file then forces that election to have a condorcet winner.
Output is in csv form and pretty printed to the screen.
