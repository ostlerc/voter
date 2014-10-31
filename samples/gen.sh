#!/usr/bin/env bash

 #peak
 for i in $(seq 1 100); do egen -cand 5 -vote 5 -peak | tally -dumb -v -o json >> peak_dumb.json; done
 for i in $(seq 1 100); do egen -cand 5 -vote 5 -peak | tally -v -o json >> peak.json; done

 #pref
 for i in $(seq 1 100); do egen -cand 5 -vote 5 -pref | tally -dumb -v -o json >> pref_dumb.json; done
 for i in $(seq 1 100); do egen -cand 5 -vote 5 -pref | tally -v -o json >> pref.json; done

 #cond
 for i in $(seq 1 100); do egen -cand 5 -vote 5 -cond | tally -dumb -v -o json >> cond_dumb.json; done
 for i in $(seq 1 100); do egen -cand 5 -vote 5 -cond | tally -v -o json >> cond.json; done

 #reg
 for i in $(seq 1 100); do egen -cand 5 -vote 5 | tally -dumb -v -o json >> reg_dumb.json; done
 for i in $(seq 1 100); do egen -cand 5 -vote 5 | tally -v -o json >> reg.json; done

 #peak pref
 for i in $(seq 1 100); do egen -cand 5 -vote 5 -peak -pref | tally -dumb -v -o json >> peak_pref_dumb.json; done
 for i in $(seq 1 100); do egen -cand 5 -vote 5 -peak -pref | tally -v -o json >> peak_pref.json; done

 #peak cond
 for i in $(seq 1 100); do egen -cand 5 -vote 5 -peak -cond | tally -dumb -v -o json >> peak_cond_dumb.json; done
 for i in $(seq 1 100); do egen -cand 5 -vote 5 -peak -cond | tally -v -o json >> peak_cond.json; done


 #results
 rm results.json
 for i in $(ls *.json | sort); do echo $i >> results.json; cat $i | poll | jq . >> results.json; done
