#!/usr/bin/env python3

# INF8775 - Analyse et conception d'algorithmes
#   TP3
#
#   Author :
#     BURLATS, Auguste - Feb 20th, 2022
#
#   Changelog:
#     28/03/2022 - Initial availability
#
#   USAGE:

#     FR - Ce script vérifie la solution qui lui est donnée pour
#          conformité avec les exigences du TP.
#          Python 3.5 ou ultérieur exigé.
#          ./tp.sh -e exemplaire -p > fichier_solution
#          ./check_sol.py -s fichier_solution -e exemplaire
#          ou 
#          python check_sol.py -s fichier_solution -e exemplaire

import sys
import re
import math
import argparse

from numpy import size


def load_instance(instance_path):
    with open(instance_path,'r') as instance_stream:
        # Process first line which defines problem characteristics
        line_one = next(instance_stream)
        dimensions = [int(i) for i in line_one.split()] #dimensions[0] : number of atomes, dimensions[1] : number of atoms types, dimensions[2] : number of edge
        next(instance_stream)
        atoms_repartition = [int(i) for i in next(instance_stream).split()]
        if(len(atoms_repartition) != dimensions[1]):
            print(1)
            return 1

        next(instance_stream)
        buffer = next(instance_stream)
        H = []
        while buffer != '\n' :
            H.append([int(i) for i in buffer.split()])
            if(len(H[-1]) != dimensions[1]):
                print(2)
                return 1
            buffer = next(instance_stream)
        if(len(H) != dimensions[1]):
            print(3)
            return 1

        liste_edge = []
        for line in instance_stream:
            if line != '\n':
                liste_edge.append([int(i) for i in line.split()])
                if liste_edge[-1][0] >= dimensions[0] or liste_edge[-1][1] >= dimensions[0]:
                    print(4)
                    return 1
                if len(liste_edge[-1]) != 2:
                    print(5)
                    return 1

        if len(liste_edge) != dimensions[2]:
            print(len(liste_edge))
            print(dimensions[2])
            return 1
    return dimensions, H, liste_edge, atoms_repartition



def is_solution_format_valid(raw_solution):
    target_pattern = r"(?:^\s*(?:\d+\s+)*\d+\s*$\n*)+"
    return bool(re.match(target_pattern, raw_solution))


def parse_solution(raw_solution):
    solution_data = []

    for line in raw_solution.splitlines():
        line_contents = line.split()

        if not line_contents:
            continue
        else:
            solution_data.append([int(x) for x in line_contents])

    return solution_data



# Error codes encapsulated into first element of tuple returned: 1 incomplete solution, 2 validity of atoms, 3 bad atoms repartition
def check_consistency(solutions, dimensions, atoms_repartition):

    for (solution, solution_index) in zip(solutions, range(len(solutions))):
        nb_atoms = [0 for i in range(len(atoms_repartition))]

        #Check the size of the line
        if len(solution) != dimensions[0]:
            return(1, solution_index)
        
        for i in solution:
            #Check if the type of i is a valid type for the problem
            if i >= dimensions[1]:
                return(2, solution_index, i)
            nb_atoms[i]+=1
        
        for index, val in enumerate(nb_atoms):
            if val != atoms_repartition[index]:
                return(3, solution_index, index, val)
        

    return 0


def compute_objective(solutions, liste_edge, H):
    solution = solutions[-1]
    objective = 0

    for edge in liste_edge:
        objective += H[solution[edge[0]]][solution[edge[1]]]

    return objective


if __name__ == '__main__':
    # Parse arguments
    parser = argparse.ArgumentParser()
    parser.add_argument("-e", "--exemplaire", \
                        help="Représente l'exemplaire correspondant à la solution étudiée", \
                        action='store', required=True, metavar='FICHIER_EXEMPLAIRE')
    parser.add_argument("-s", "--solution", \
                        help="Représente la solution", \
                        action='store', required=True, metavar='FICHIER_SOLUTION')
    args = parser.parse_args()

    # Load instance corresponding to solution
    instance_data = None
    try:
        instance_data = load_instance(args.exemplaire)
    except:
        print("Erreur : impossible d'ouvrir le fichier de l'exemplaire.", file=sys.stderr)
        sys.exit(1)
    
    if instance_data == 1 and instance_data is not None:
        print("Erreur : l'exemplaire fourni en argument à ce script de vérification a un format non valide. "\
            "Vérifiez le chemin et/ou le contenu de l'exemplaire.", file=sys.stderr)
        sys.exit(1)

    with open(args.solution, 'r') as fichier:
        solution_content = fichier.read() 

    # Check whether format is as expected 
    if not is_solution_format_valid(solution_content):
        print("Erreur : les solutions fournies en pipe à stdin ont un format non valide. Revoyez la convention discutée dans l'énoncé.", file=sys.stderr)
        print("A reçu :", file=sys.stderr)
        print(solution_content, file=sys.stderr)
        sys.exit(1)


    # Structure piped solution in memory
    resolution_data = parse_solution(solution_content)
    
    # Check solutions' consistency
    consistency_result = check_consistency(resolution_data, instance_data[0], instance_data[3])
    if consistency_result != 0:
        print("Erreur : une ou plusieurs des solutions fournies en pipe à stdin présentent un problème de consistance.", file=sys.stderr)

        if consistency_result[0] == 1:
            print("Raison : la solution " + str(consistency_result[1]) \
                  + " (0-indexé) contient un nombre inadéquat d'atomes.", file=sys.stderr)
        elif consistency_result[0] == 2:
            print("Raison : l'atome de type " + str(consistency_result[2]) + " de la solution " + str(consistency_result[1]) \
                  + " (0-indexés) n'est pas d'un type valide. Les types acceptés pour cette instance sont de 0 à " + str(len(instance_data[3]) - 1) \
                      +" compris.", file=sys.stderr)
        elif consistency_result[0] == 3:
            print("Raison : la solution " + str(consistency_result[1]) + " (0-indexés) ne contient pas le bon nombre d'atomes du type "  \
                  + str(consistency_result[2]) + ". Nombre attendu : " + str(instance_data[3][consistency_result[2]]) \
                  + ". Nombre présent : " + str(consistency_result[3]), file=sys.stderr)

        sys.exit(1)

    # Satisfied by the solutions' presentation, compute best objective
    objective = compute_objective(resolution_data, instance_data[2], instance_data[1])
    print("OK : la valeur de l'objectif de la dernière (ie, meilleure) solution fournie est de " + str(objective) + ".\n")
