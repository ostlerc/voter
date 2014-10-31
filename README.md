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
      -o="csv": output type [json,csv]
      -t="borda,bucklin,copeland,slater,kemeny,stv": tally type results
      -v=false: verbose output. Show all tally information

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

    - sample (verbose and pretty printed) csv output
        weight          5      4      3      6
        Alex            5      1      6      4
        Bart            1      6      5      5
        Cindy           2      3      7      3
        David           4      4      1      2
        Erik            6      5      3      1
        Frank           3      2      2      6
        Greg            7      7      4      7

        pref            5      2
        rank            borda  bucklin  copeland  kemeny  slater  stv
        manipulation    true   true     true      true    true    true
        pref intact     true   true     true      true    true    true
        irrlvnt alters  false  false    false     false   false   false
        1               Erik   Erik     Alex      Erik    Alex    Erik
        2               Alex   Bart     -1        Alex    David   Alex
        3               Frank  -1       -1        David   Erik    David
        4               David  -1       -1        Frank   Bart    Frank
        5               Cindy  -1       -1        Cindy   Frank   Bart
        6               Bart   -1       -1        Bart    Cindy   Cindy
        7               Greg   -1       -1        Greg    Greg    Greg

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

poll
====
  poll verbose tally outputs to find trends

    stdin = newline separated stream of verbose tally json

    - sample json output
    {
        "total": 100,
        "condorcet_lost": {
            "stv": 25,
            "slater": 5,
            "kemeny": 3,
            "bucklin": 16,
            "borda": 25
        },
        "irr_cand_affect": {
            "stv": 98,
            "slater": 98,
            "kemeny": 98,
            "copeland": 98,
            "bucklin": 98,
            "borda": 98
        },
        "manipulations": {
            "slater": 95,
            "kemeny": 99,
            "copeland": 95,
            "bucklin": 93,
            "borda": 99
        },
        "pref_changed": {
            "stv": 4
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

    cat sample.csv | tally -i csv | column -s, -tn

In this example we tally an election from a csv file and output it in csv format. Then pretty print it.

    egen | tally -o json | jq .

In this example we generate a random election and tally it into a json result. The result is then pretty printed

    cat sample.csv | egen -cond -o csv | column -s, -tn

This example generates an election from a csv file then forces that election to have a condorcet winner.
Output is in csv form and pretty printed to the screen.

    for i in $(seq 1 100); do egen -vote 5 -cand 5 | tally -v -o json >> res.json; done; cat res.json | poll | jq .

In this example we create 100 elections and verbose tally them. The elections are stored in a file called 'res.json'.
After the elections are created we then pipe them into poll and receive our pretty printed json results

Generation
==========

For some sample ways to generate poll files, you can look at samples/gen.sh. The samples folder contains all the
poll files as examples for using this software.

Results
=======

I have generated a fair amount of elections and placed them inside of the samples folder. For each of these
elections I did a verbose tally and printed the results in the samples/results.json file. I have noticed that
manipulating an election is much easier than I had anticipated. In fact, I was so surprised that I generated a
few different ways of scoring (comparing) votes to see if a dumb way would provide different results. My first
method of scoring is complicated and smart. Weighting top candidates and lower candidates out of position
negatively. The next voting mechanism I made is called 'dumb'. You can provide the dumb flag to the tally
binary and it will use the dumb scoring method instead of the smart one. The results were shockingly similar.
The dumb scoring method only says a manipulation is successful if your top candidate was not in first, and by 
voting differently your top candidate has now become the winner. As you can see, manipulations are around 95%
possible by any of the randomly generated elections I have seen. The exception to this is the stv, when scoring
scoring with dumb then little to no manipulations were found. This was especially true when condorcet winners 
or peak voting was occuring.


I played around with candidate size and vote size but it seemed to have minor effects on the numbers.
Generating large elections is easy, but tallying and find manipulations is a slow processes as most of The
algorithms are polynomial and poorly written. (by me)

The changing of preferences is more difficult than manipulation. That seems to only be around 10% of the time
possible. Condorcet winners can also lose quite a lot - up to 50% of my elections had a voting scheme that
allowed the condorcet winner to not be the winner. Kemeny seems to be the most resilient to manipulations.
Removing irrelevant candidates changed the top winner almost 90% of the time. Borda has shown to be more
difficult to make irrelevant candidates change the top winner, but it still has around 45% chance of that
affecting the results.

Preferences changing was very rare. STV seemed to be the most likely scheme where this will happen, and even
so it is only around 13% of the time. Copeland has around a 5% chance that preferences could not be preserved.
