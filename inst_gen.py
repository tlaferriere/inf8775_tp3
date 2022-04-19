#!/usr/bin/env python3

# INF8775 - Analyse et conception d'algorithmes
#   TP3 - Configuration d'atomes
#
#   AUTEUR :
#     PEZZOLI, Gauthier - 26 mars 2022
#
#
#   USAGE :
#     Ce script génère les exemplaires requis pour le TP3.
#
#     $ ./inst_gen.py [-h] -t NB_SITES -k NB_TYPES [-n NB_EXEMPLAIRES]
#
#     où :
#       * NB_SITES est le nombre de site (également le nombre d'atomes)
#       * NB_TYPES est le nombre de type d'atomes
#       * NB_EXEMPLAIRES est le nombre d'exemplaires différents requis (par défaut 1).
#
#     Il est nécessaire de rendre ce script exécutable en utilisant chmod +x
#     Python 3.5 ou ultérieur recommandé pour lancer ce script.

import random
import argparse
import numpy as np


if __name__ == "__main__":
    # Parse arguments
    parser = argparse.ArgumentParser()
    parser.add_argument("-t", "--taille", \
                        help="Représente le nombre de site (également le nombre d'atomes)", \
                        action='store', required=True, metavar='NB_SITES', type=int)
    parser.add_argument("-k", "--types", \
                        help="Représente le nombre de type d'atomes", \
                        action='store', required=True, metavar='NB_TYPES', type=int)
    parser.add_argument("-n", "--nb-exemplaires", \
                        help="Représente le nombre d'exemplaires d'une même taille à générer", \
                        action='store', required=False, metavar='NB_EXEMPLAIRES', type=int)

    args = parser.parse_args()
    if not args.nb_exemplaires:
        args.nb_exemplaires = 1

    if args.types > 6 or args.types < 2:
        raise NameError('k doit etre entre 2 et 6')
    
    
    for num in range(args.nb_exemplaires):
        densite = 0.2
        NB_ARRETES = 0
        # Generation matrice adjacence
        adj = np.zeros((args.taille,args.taille))
        for i in range(args.taille):
            for j in range(i+1,args.taille):
                if random.random() < densite:
                    adj[i,j] = 1
                    adj[j,i] = 1
                    NB_ARRETES +=1


        connex = np.zeros((args.taille))
        sec = 1
        comp = []
        pile = [0]
        connex[0] = 1
        while True:
            while pile:
                ele = pile.pop(0)
                comp.append(ele)
                for j in range(args.taille):
                    if (adj[ele,j] and connex[j] == 0):
                        connex[j] = 1
                        pile.append(j)

            if len(comp) < args.taille:
                for autre in range(args.taille):
                    if connex[autre] == 0:
                        connex[autre] = 1
                        sec = autre
                        break
                pre = comp[random.randint(0,len(comp)-1)]
                if adj[pre,sec]:
                    print('error')
                adj[pre,sec] = 1
                adj[sec,pre] = 1
                NB_ARRETES += 1
                pile.append(sec)
            else:
                break
                
        # K Types
        K_type = [0]*args.types
        for _ in range(args.taille):
            K_type[random.randint(0,args.types-1)] += 1
        # Mélange
        n_exchange = random.randint(0,10)
        for _ in range(n_exchange):
            a = random.randint(0,args.types-1)
            b = random.randint(0,args.types-1)
            while b == a:
                b = random.randint(0,args.types-1)
            exchange = K_type[a]//2
            K_type[a] -= exchange
            K_type[b] += exchange

        # H
        H = np.zeros((args.types,args.types))
        for i in range(args.types):
            for j in range(i,args.types):
                H[i,j] = int(round(np.random.normal(2,5),0))
                H[j,i] = H[i,j]
        
        
        with open('N' + str(args.taille) + '_K' + str(args.types) + '_' + str(num),'w') as inst:
            inst.write(f'{args.taille} {args.types} {NB_ARRETES}\n')
            inst.write('\n')
            inst.write(' '.join([str(k) for k in K_type])+'\n')
            inst.write('\n')
            for i in range(args.types):
                inst.write(' '.join([str(int(H[i,j])) for j in range(args.types)])+'\n')
            inst.write('\n')
            for i in range(args.taille):
                for j in range(i+1,args.taille):
                    if adj[i,j] == 1:
                        inst.write(str(i) + ' ' + str(j)+'\n')

                
