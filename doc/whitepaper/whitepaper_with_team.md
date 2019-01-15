---
abstract: |
  Data is the new oil. Every day we create 2.5 quintillion bytes of data[@keyIBM].
  90 percent of the data in the world today has been created in
  the last two years alone – and with new devices, sensors
  and technologies emerging, the data growth rate will likely accelerate
  even more.

  Not only centralized SW giants, but also mobile and network operators and
  various enterprises that install huge number of devices
  or any electronic infrastructure are in position to put the sensors in their
  equipment and collect huge amount of data. Some of this data is perishable -
  i.e. it must be consumed instantly or it looses value. Some of this data is long-lasting. No mater what kind of data stakeholders collect, they usually have the same problem:
  how to draw additional profit from this data, beyond it's immediate and obvious purpose[^1].

  All this data would have value for many parties and can be further monetized.
  Data collectors could become data sellers, and offer collected data on the
  specialized data marketplace. On the other hand, data buyers would be interested
  to browse offered data streams and buy them, then use this data to further
  process it and/or build new services for their customers.

  **An eBay for IoT sensor data is need. This eBay is called Datapace**.

  [^1]: Take for example a mobile telephony operator. Company like this already owns
  huge number of network base-stations, gateways and antennas which make the
  deployed network infrastructure. These network devices are already equipped
  with big number of telemetry sensors that provide the operating state insights
  and are used for management and maintenance. Data coming in for this sensors
  is useful for the operator to keep the network healthy and functional.
  But beyond this primary purpose, collected data can be extremely useful
  for other parties - like Smart City municipalities, health institutions
  or various other businesses. Moreover, because of the density of mobile
  base-stations and antennas, operators are in unique position to offer for example
  extremely precise environmental data, which is hard to achieve to even specialized
  services as it demands significantly expensive HW sensor installation. Similarly, a company that does smart signage could turn public signs into
  sensing stations with marginal additional costs. With further cost drop of the sensors and appearance of smart dust, even individuals or small enterprises will be capable to collect significant amount of IoT data.
---

# Introduction
Datapace is an eBay for IoT sensor data. But beyond IoT, Datapace marketplace
can be used to sell or buy any type of data, independently of it's type or provenance.

Datapace is distributed and decentralized system based on blockchain.
Blockchain technology is used for several important purposes in Datapace system:

- To enable tokenization of value (i.e. provide TAS token) and token economy
- To insure data integrity (i.e. to store data hashes and guarantee that data is not tampered with)
- To enable Smart Contract capabilities
- To provide network security via PBFT consensus and immutability and
make the system hack-proof

Each of these characteristics of blockchain and how they are leveraged upon in the Datapace system will be explained in more details in the following chapters.

Datapace market place is built with intention to be simple, easy to use and intuitive.
Anyone familiar with eBay should immediately understand how to sell digital assets -
in this case the data stream, or how to browse offered data streams and purchase selected data. Simplicity of use opens possibility for mass-market adoption while simplicity of the system provides high quality implementation and better secured and more performant application.

# Stakeholders
Datapace system is built on private, permissioned blockchain. It uses PBFT algorithm
for concensus and state replication, which guarantees high transaction throughput and
fast transaction finality (which as a consequence prevents blockchain forking). Because of the nature of PBFT algorithm, the whole system is run by a closed consortium with a known set of validators. Never the less, any entity can potentially request access to the consortium and run a validating node under contractual agreements.

Based on this we can identify following stakeholders of Datapace system:

- Data buyers
- Data sellers
- Validators

## Data Buyers
Data buyers are organizations and individuals that are interested in buying the data. They log into the system and browse the data streams offered for sell, as one would browse items on eBay or Amazon marketplace.

Data streams are offered under certain price and can are purchased for TAS tokens.

Data buyer must have a sufficient amount of TAS tokens in his wallet in order to purchase the data. Once data is purchased, data buyer obtains a proxied HTTP URL from which he can consume the data. This URL is unique and temporary - it expires after the lease period for which data was payed for.

## Data Sellers
Data sellers are organizations or individuals that offer the data for sell.

It is responsibility of data seller to provide a valid data source URL and give detailed description of the data stream an it's format (it's JSON schema) - so it can be easily consumed by data buyer. This URL is secret, and it is never reveled to data buyer. It is only temporary proxy URL that is given to data buyer, and it expires after time data was payed for.

Additionally, data seller can provide geolocation data of the stream source, so that it can be queried on the maps.

Data sellers should provide valid data sources. In order to guarantee the validity of the data, Datapace employs several mechanisms - like seller reputation rating and verified IoT gateway HW provisioning, which will be explained in dedicated chapter.

Data sellers obtain TAS tokens in their wallet when the stream that they offered is purchased.

## Validators
Validators are the entities that participate in network infrastructure, i.e. in block validation. Validators are rewarded for their work in TAS tokens.

Because in the phase 1 Datapace is based on private PBFT blockchain, set of validators must be known up-front. Datapace consortium will allow adherence of new members under strict contractual agreements.

In the second phase of development, Datapace validation will be opened to public via novel _Proof-of-Verified-Source_ and _Proof-of-Stake_ on the Cosmos[@cos] network.

# System Architecture
## Description
Datapace is a decentralized application based on the blockchain network with native token of value.

Datapace blockchain is provided via novel BigchainDB technology, which provides several benefits to the system:

- Native token, in a form of divisible digital asset
- Digital asset queries
- Fast transaction throughput and finalization

\begin{figure}
\begin{center}
\begin{tikzpicture}[>=stealth]

  %nodes
  \node[draw, minimum width=3cm, minimum height=2cm, anchor=center, text width=2cm, align=center, fill=gray!20] (BACKEND) {Datapace\\Backend};

  \node[draw, minimum width=3cm, minimum height=2cm, xshift=-2.5cm, yshift=3cm, left of=BACKEND, text width=2cm, align=center, fill=gray!20] (UI) {Datapace\\UI};

  \node[draw, minimum width=3cm, minimum height=2cm, yshift=2cm, above of=UI, text width=2cm, align=center] (BROWSER) {Browser};

  \node[draw, minimum width=3cm, minimum height=2cm, xshift=2.5cm, yshift=3cm, right of=BACKEND, text width=2cm, align=center, fill=gray!20] (PROXY) {Datapace\\Proxy};

  \node[draw, minimum width=3cm, minimum height=2cm, yshift=2cm, above of=BACKEND, text width=2cm, align=center, fill=gray!20] (MFX) {Datapace\\IoT\\Platform};

  \node[draw, minimum width=2cm, minimum height=2cm, xshift=-1.5cm, yshift=2cm, above of=PROXY, text width=2cm, align=center] (DS) {Data\\Source};

  \node[draw, minimum width=2cm, minimum height=2cm, xshift=1.5cm, yshift=2cm, above of=PROXY, text width=2cm, align=center] (DR) {Data\\Consumer};

  \node[draw, minimum width=5cm, minimum height=2cm, yshift=-1.5cm, below of=BACKEND, text width=2cm, align=center, fill=gray!20] (BC) {Blockchain};

  % draw the paths and and print some Text below/above the graph
  \path (BACKEND) edge[bend left=40] (UI);
  \path (BACKEND) edge[bend right=40] (PROXY);
  \path (BACKEND) edge[-] (MFX);
  \path (MFX) edge[-] (DS);
  \path (PROXY) edge[-] (DS);
  \path (PROXY) edge[-] (DR);
  \path (UI) edge[-] (BROWSER);
  \path (BACKEND) edge[-] (BC);
\end{tikzpicture}
\end{center}
\caption{Datapace System Architecture}
\label{fig:arch}
\end{figure}


In order to secure the system and make it resistant to Byzantine General attack[@bft], Datapace replaces the underlying BigchainDB default BCA consensus algorithm[@bcdb] with Tendermint PBFT consensus engine[@tendr].

Tendermint is very performant PBFT consensus algorithm - it supports thousands of transaction per second at 1000ms latencies. Not only that Datapace benefits from this consensus algorithm in security and performance, but this mechanism opens the possibility to connect Datapace system to incoming Cosmos network. Announced as "Internet of Blockchains", Cosmos hub will give to Datapace system two very important features: interoperability and additional scalability.

\begin{figure}
\begin{center}
\begin{tikzpicture}[>=stealth]

  %nodes
  \node[draw, minimum width=3cm, minimum height=1cm, anchor=center, text width=3cm, align=center] (BDB1) {BigchainDB\_1};
  \node[draw, minimum width=3cm, minimum height=1cm, below of=BDB1, text width=3cm, align=center, fill=gray!20] (T1) {Tendermint\_1};

  \node[draw, minimum width=3cm, minimum height=1cm, anchor=center, xshift=-4cm, yshift=-3cm, text width=3cm, align=center] (BDB2) {BigchainDB\_2};
  \node[draw, minimum width=3cm, minimum height=1cm, below of=BDB2, text width=3cm, align=center, fill=gray!20] (T2) {Tendermint\_2};


  \node[draw, minimum width=3cm, minimum height=1cm, anchor=center, xshift=4cm, yshift=-3cm,text width=3cm, align=center] (BDB3) {BigchainDB\_3};
  \node[draw, minimum width=3cm, minimum height=1cm, below of=BDB3, text width=3cm, align=center, fill=gray!20] (T3) {Tendermint\_3};

  \node[draw, minimum width=3cm, minimum height=1cm, anchor=center, xshift=-2.5cm, yshift=-6cm, text width=3cm, align=center] (BDB4) {BigchainDB\_4};
  \node[draw, minimum width=3cm, minimum height=1cm, below of=BDB4, text width=3cm, align=center, fill=gray!20] (T4) {Tendermint\_4};


  \node[draw, minimum width=3cm, minimum height=1cm, anchor=center, xshift=2.5cm, yshift=-6cm,text width=3cm, align=center] (BDB5) {BigchainDB\_5};
  \node[draw, minimum width=3cm, minimum height=1cm, below of=BDB5, text width=3cm, align=center, fill=gray!20] (T5) {Tendermint\_5};

  % draw the paths and and print some Text below/above the graph
  \path (T1) edge[bend right=43] (BDB2);
  \path (T1) edge[bend left=43] (BDB3);

  \path (T2) edge[bend right=60] (T4);
  \path (T3) edge[bend left=60] (T5);

  \path (T4) edge[-] (T5);

  \path (T2) edge[-] (BDB5);
  \path (T3) edge[-] (BDB4);

  \path (T1) edge[-] (BDB5);
  \path (T1) edge[-] (BDB4);

  \path (BDB2) edge[-] (BDB3);

\end{tikzpicture}
\end{center}
\caption{Datapace Blockchain Network}
\label{fig:arch}
\end{figure}

Interoperability is extremely important, as it will enable TAS token to natively flow from Datapace private blockchain into other blockchains connected to the Cosmos hub, thus opening potential for TAS exchange to other crypto-currencies, and vice versa. This will influence token economy and raise the value of the TAS token. Additionally, developed token economy would allow _Proof-of-Stake_ consensus to be applied on the top of the Datapace-Tendermint system and allow opening Datapace validator set participation to the wide public.

Scalability is also important, although, as a consequence of the wise technology choices, Datapace system is already very performant. But "Interent of Blockchains" will enable additional scaling od Datapace chains through sharding[@shard] using Cosmos zones.

\begin{figure}
\begin{center}
\begin{tikzpicture}[>=stealth]

  %nodes
  \node[draw, minimum width=3cm, minimum height=1cm, anchor=center, text width=2cm, align=center] (M1) {Datapace Blockchain 1};

  \node[draw, minimum width=3cm, minimum height=1cm, anchor=center, text width=2cm, align=center, below of=M1, yshift=-1cm] (M2) {Datapace Blockchain 2};

  \node[draw, minimum width=3cm, minimum height=1cm, anchor=center, text width=2cm, align=center, below of=M2, yshift=-1cm] (M3) {Datapace Blockchain 3};

  \node[draw, minimum width=2cm, minimum height=1cm, anchor=center, text width=2cm, align=center, right of=M2, xshift=3.5cm, circle] (COS) {Cosmos Hub};

  \node[draw, minimum width=3cm, minimum height=1cm, anchor=center, text width=2cm, align=center, right of=COS, xshift=3.5cm, yshift=2cm] (B) {Bitcoin Blockchain};

  \node[draw, minimum width=3cm, minimum height=1cm, anchor=center, text width=2cm, align=center, below of=B, yshift=-1cm] (E) {Ethereum Blockchain};

  \node[draw, minimum width=3cm, minimum height=1cm, anchor=center, text width=2cm, align=center, below of=E, yshift=-1cm] (C) {Custom Blockchain};

  % draw the paths and and print some Text below/above the graph
  \path (COS) edge[-] node[anchor=north,above]{TAS} (M1);
  \path (COS) edge[-] node[anchor=north,above]{TAS} (M2);
  \path (COS) edge[-] node[anchor=north,above]{TAS} (M3);

  \path (COS) edge[-] node[anchor=north,above]{BTC} (B);
  \path (COS) edge[-] node[anchor=north,above]{ETH} (E);
  \path (COS) edge[-] node[anchor=north,above]{XXX} (C);

\end{tikzpicture}
\end{center}
\caption{Datapace Sharding and Interoperability via Cosmos}
\label{fig:arch}
\end{figure}


As mentioned before, blockchain technology is used for several important purposes in Datapace system:

- **TAS token**: TAS token is native token of value in Datapace system and is necessary for system operation and functioning. It will be explained in details in a dedicated chapter.
- **Data integrity**: Leveraging on BigchainDB digital asset features, as well as native digital asset querying, Datapace implements mechanism that insures integrity of the data that flows through the system by taking it's digital fingerprint (cryptographic hash) and stores it in to the immutable blockchain database. This way system assures that critical data has not been tampered with. In the context of OTA firmware updates of safety-critical IoT devices or tamper-proof checking of already running software on such a systems (for example a braking system of a self-driving vehicule) this form of data security becomes quintessential.
- **Smart Contract**: Smart Contracts define a complex set of conditions under which data is exchanged. They are important part of Datapace system, and will be explained in detail in a dedicated chapter.
- **Network security (via PBFT consensus)**:  In order to protect valuable digital assets and network infrastructure in the era of ever-increasing security threats[^2], Datapace builds a decentralized network based on Byzantine fault-tolerant state and data replication algorithm. This way system can tolerate up to 1/3 malicious-acting nodes and assure network functioning under cyber-attack. Additionally, blockchain-structured data assure immutability and anti-tampering characteristics. Applying _Proof-of-Validated-Source_ and _Proof-of-Stake_ consensus, network is adding an additional layer of protection, incentivizing nodes to behave honestly and punishing badly behaving nodes. Based on these important features and technologies, Datapace builds high-security network that is capable to fully protect digital assets and insure secure protection of value exchanged through Datapace marketplace.
- **Auditing (via record immutability)**: Datapace enables monetary transactions, which are often subject to various regulations and can be examined by regulatory bodies. Thanks to the immutability feature of blockchain systems, Datapace system allows every organization participating in Datapace data market to have a proven track of records of all executed transactions.

To make system usable for the wide public, Datapace implements secure centralized wallet, similar to Coinbase[@cbase]. Wallet, however, can be implemented also in decentralized fashion, as users can create accounts and transfer their funds to wallet of their choice.

In order to facilitate and standardize IoT data collection, Datapace system provides an IoT platform, based on Mainflux[^mfx]. Mainflux IoT platform is integrated into Datapace system part called "Datapace IoT Platform", and exposes an API for sensor connection and management. Mainflux (and thus Moentasa IoT platform) equally provides IoT messaging and persistence capabilities, so all the data from sensors can be either offered in real-time or stored for historical usage.

Users of Datapace can use these interfaces to easily connect their sensors and thus enable data collection - this way they do not have to go through IoT system set-up, but can use Datapace platform as a service.

Additionally, once data enters the Datapace system via IoT platform, various types of processing and data analytics can be applied. Datapace will offer AI and ML algorithms to be applied to collected IoT data, so that users can enrich their data with with different type of intelligence prior to offering it to the market. This will boost the price of their data streams.

Finally, data coming through Datapace IoT platform is in standardized format - as per Mainflux specification, data is formatted in SenML[^senml] - a media type definition for representing simple sensor measurements and device parameters which comes in JSON and CBOR[^cbor] flavors. This facilitates the operations for data consumers - they know up front what data format to expect and the same set of processing scripts, procedures and programs can be applied to various data streams.

It is important to note that use of Datapace IoT platform is not mandatory, even though it brings obvious benefits. Users can still obtain their own data via their legacy IoT installations and provide only access link to this data via Datapace market.

[^2]: The number of records compromised grew a historic 566 percent in 2016 from 600 million to more than 4 billion. These leaked records include data cybercriminals have traditionally targeted like credit cards, passwords and personal health information, but IBM study[@ibmSaf] also shows a shift in cybercriminal strategies. In 2016, a number of significant breaches related to unstructured data such as email archives, business documents, intellectual property and source code were also compromised.

[^mfx]: Mainflux[@mfx] (<https://www.mainflux.com>) is modern, scalable, secure open source and patent-free IoT cloud platform written in Go. It accepts user, device, and application connections over various network protocols (i.e. HTTP, MQTT, WebSocket, CoAP), thus making a seamless bridge between them. It is used as the IoT middleware for building complex IoT solutions.

[^senml]: SenML[@senml] is a sensor markup language that aims to simplify gathering data from different devices across the network. It simply is JSON containing named events together with an associated value and unit.

[^cbor]: CBOR (Concise Binary Object Representation) is a binary data serialization format loosely based on JSON. It is defined in IETF RFC 7049[@cbor]

## Technology Summary
Datapace integrates several open-source technologies which in combination provide a powerful system. An overview of technologies used is given in the table \ref{tab:table1}.

\begin{table}[h!]
\caption{Datapace technology summary}
\label{tab:table1}

\begin{center}

\begin{tabularx}{\textwidth}{>{\bfseries}lX}

\toprule

BigchainDB        & Blockchain (Distributed Ledger). Provides token as a form of divisible digital asset. Immutability, querying, validator voting. Fast transactions.    \\
\midrule

Tendermint        & Provides PBFT consensus algorithm and P2P machine state replication. Adds security to Datapace blockchain (BigchainDB network). Facilitates Cosmos integration.     \\
\midrule

Cosmos            & Provides TAS token interchangeability. Provides Datapace blockchain scalability through sharding. Provides interoperability with other blockchain networks - like Ethereum or Bitcoin.      \\
\midrule

Mainflux          & Provides IoT platform as a service. Enables IoT sensor and gateway connectivity and management. Provides IoT messaging, real-time and persisted data. \\

\bottomrule

\end{tabularx}
\end{center}
\end{table}

# Data Verification
## Overview
Datapace has unique-on-the-market solution for verifying the source of the IoT data. Based on the fact that Datapace and it's partners play one of the crucial roles in telecom equipment industry, an IoT gateways and edge computers were designed and connected with big number of sensors to serve as a verified and known IoT sensor data source.

Datapace installs these sensors in cooperation with network and telecom partners, or sends the certified equipment to various other partners for installation. Because these edge computers, IoT gateways and sensors contain known and certified hardware and firmware, often coupled with embedded GPS modules, system can be assured that data coming from these sensors is:

- Real-world data and not modified or generated "fake" data
- Coming from precise geographical location

Datapace partners that install and deploy this equipment will have an advantage on the marketplace, as their data sources will be marked as "trusted and verified".

Moreover, since these partners made an economic investment and also entered in partnership with Datapace through various legal contractual agreements, they are allowed to run a validating node and participate in _Proof-of-Verified-Source_ network consensus. Validators are rewarded for their work with TAS tokens.

## Implementation
Datapace sensors are attached to Datapace gateways and edge computers, or directly connected to the IoT platform.

Datapace provides IoT platform to monitor and manage IoT network and gather the data from the installed sensors. Once data is collected, it can be offered for sell by the user of Datapace system that has installed sensors (and/or gateways) and is the owner of the data.

\begin{figure}
\begin{center}
\begin{tikzpicture}[>=stealth]

  %nodes
  \node[draw, minimum width=3cm, minimum height=3cm, anchor=center, text width=2cm, align=center] (MFX) {Datapace\\IoT\\Platform};

  \node[draw, minimum width=1cm, minimum height=1cm, anchor=center, text width=1cm, align=center, left of=MFX, xshift=-5cm, yshift=1cm] (S1) {Sensor 1};

  \node[draw, minimum width=1cm, minimum height=1cm, anchor=center, text width=1cm, align=center, below of=S1, yshift=-1cm] (S2) {Sensor 2};

  \node[draw, minimum width=2cm, minimum height=1cm, anchor=center, text width=2cm, align=center, right of=S1, xshift=1cm] (GW) {Datapace Gateway};

  \node[draw, minimum width=3cm, minimum height=3cm, anchor=center, text width=2cm, align=center, right of=MFX, xshift=3cm] (MON) {Datapace Marketplace};


  % draw the paths and and print some Text below/above the graph
  \path (S1) edge[-] (GW);
  \path (GW) edge[-] (MFX);
  \path (S2) edge[-] (MFX);
  \path (MFX) edge[-] (MON);

\end{tikzpicture}
\end{center}
\caption{Datapace IoT device management via IoT platform}
\label{fig:arch}
\end{figure}

An additional benefit of enabling data connection via Datapace IoT platform is that users that choose this option can add various data processing and analytic services offered by Datapace system. Additionally, they can apply ML and AI insights to their data prior to offering them for sell, which would significantly augment the data price.

# TAS Token
TAS token is utility token of Datapace system. It is used to assure fair and secure functioning of the system, as well as to enable token economy on the Datapace data market.

Primary purpose of the token will be to fuel the system - it will be used to tokenize the value of digital assets (i.e. data) and facilitate their exchange. Equally, once the token economy is developed, TAS token will have a purpose in enabling the consensus mechanism based on _Proof-of-Stake_.

Data sellers will use TAS token as a representation of value of their digital data streams that are offered on the marketplace. Buyers will use TAS token to exchange it for selected data - they will transfer their TAS tokens to data sellers and obtain their digital assets in return.

# Proof-of-Verified-Source and Proof-of-Stake
Datapace system employs two types of proof schemes that allow data providers and network participants to prove the data origin and quality as well as to enforce
honest behavior of data sellers.

## Proof-of-Verified-Source
In order to secure the network, Datapace system provides and original an unique approach called _Proof-of-Verified-Source_. This approach represents validator (miner) selection algorithm based on a proof of monetary investment in sensing hardware and networking equipment.

Due to the unique position on the market Datapace produces and delivers to the companies and network operators a specialized networking and sensing equipment - often an IoT edge gateway inter-connected with a lot of sensors. Since this hardware (and internal secure firmware) comes from known source (Datapace company), and since all equipment purchase and installation is done according to contractual agreements, everybody can be assured that this given data source is valid.

Because a company or an operator that purchases the equipment has to invest money, and also respect the written legal contracts, system can stay assured with high probability that they are incentivised to make fair decision (it is in their best of interest to keep the network secure and functional - otherwise their investment will be useless and they will suffer legal penalties).

Moreover, possibility to have _Verified Data Source_ badge listed next to the data sources offered by these companies is an additional incentive for them to purchase the specialized sensors and other equipment.

## Proof-of-Stake
Once TAS token economy is developed, a _Proof-of-Stake_ consensus algorithm will be applied in order to additionally incentivise companies and individuals that run validator nodes and help secure the network.

_Proof-of-Stake_ will equaly be used to enforce honest behavior of data sellers,
because they will have to invest a monetary deposit (in form of TAS tokens).
In case of malicious behavior (wrong data delivered, or data not delivered at all)
deposit will be withdrawn by the system and bad actor will be punished.

# Smart Contracts
A Smart Contract is a computer protocol intended to facilitate, verify, or enforce the negotiation or performance of a contract.  The aim with Smart Contracts is to provide security that is superior to traditional contract law and to reduce other transaction costs associated with contracting.

Datapace platform provides possibility for users to define and deploy Smart Contracts that automate processes and formalize contractual agreements regarding various features of the system. One important feature, for example, is revenue sharing - every data seller can define a Smart Contract that will be signed by his partners and himself. Earnings obtained by selling this data stream will then be automatically divided between the parties, without further intervention from the seller and his partners.

Moreover, Smart Contracts enable fine-grained per-user and per-datastream conditions to be formalized. For example, new GDPR (_General Data Protection Regulation_)[@gdpr] laws by which by which the European Parliament, the Council of the European Union and the European Commission intend to strengthen and unify data protection of individuals, regulate the way that telecom operators or other companies can share user data. Since this data and it's sharing and monetization represent a core business of many companies (especially of those who's business model is based on advertising), a strict new relation between company and it's users is imposed and can be formalized and automated via Smart Contracts.

Datapace UI will enable defining these Smart Contracts in a simple manner though well-defined forms. Moreover, Datapace API will provide possibilities for these contracts to be defined and deployed programatically.

# Data Integrity Through Anchoring
It is very well known feature of blockchains to offer immutable data storage. Once data is written in the blockchain it can not be changed (tampered with). This feature can be used to prove integrity of the data, which is especially important for OTA (_Over-the-Air_) firmware updates of safety-critical IoT devices or tamper-proof checking of already running software on various robots, machines, vehicles and similar.

In order to enable this feature as a service, Datapace implements an API on the top of its system that allows "anchoring" the data timestamp and cryptographic hash into the blockchain. This cryptographic hash essentially represents digital fingerprint of the data. Data hash can be recalculated and compared to immutable record in the blockchain at any later point, thus proving that the data has not be tampered with.

# Future Work - Computing and Storage Tokenization
Besides data, a marketplace based on the blockchain can allow economy of at least two important resources:

- Storage
- Comuputing

## Storage
Companies like **Storj**[@storj] or **Sia**[@sia] announced projects that strive to enable decentralized cloud. With low prices that would be a consequence of tokenized storage capacity offered by the various users in exchange of tokens, these companies can become a real competitors of SW giants in the cloud bussiness space - like Amazon or Google.

Datapace plans to integrate and maintain permissioned distributed file-system through wich Datapace users will be capable to offer and rent their storage space in exchange for TAS tokens.

## Computing
Projects like **Golem**[@golem] or **SONM**[@sonm] are working on decentralizing the computing power.

Based this ideas, but also on the ideas presented by **Blue Horizon** project from IBM [@bhz], Datapace plans to enable Docker container based decentralized platform for deploying arbitrary software on the computing infrastructure offered and rented by Datapace users in exchange of TAS tokens.

# Conclusion
Based on many reports[^mck], we can be sure of one thing: there is gold in the mountains of data. A way is needed to mine all this gold - a platform is needed to monetize all this data. Datapace is a an enabler that will unlock this huge potential.

Datapace builds decentralized marketplace based on blockchain, that is secure and scalable. It enables new token economy - TAS token will be used as an utility token of Datapace system, and will be used to enable fair and secure functioning of the system as well to enable trading facilities.

Datapace builds whole environment needed for quick adoption of the system: UI, wallet, API and SDKs. This will lower adoption barriers and lead to the higher popularity of the system, which will in turn incentivise the economy based on TAS token.

Due to unique positioning, Datapace provides specialized senor hardware, and employing various patent-pending techniques assures that data sources are verified. Moreover, through specific AI and machine learning algorithms, Datapace system assures that all data streams can be unified in format and prepared for easy consummation. This brings clear advantage of Datapace comparing to all existing competition.

Datapace will be go-to marketplace for data monetization - any data, anywhere.

[^mck]: IDC says that worldwide revenues for big data and business analytics will grow from $130.1 billion in 2016 to more than $203 billion in 2020, at a compound annual growth rate (CAGR) of 11.7%[@idc]. In addition to being the industry with the largest investment in big data and business analytics solutions (nearly $17 billion in 2016), banking will see the fastest spending growth. New report from McKinsey & Company's Global Institute is trying to put a real dollar amount to the global IoT market. In the report's estimation, IoT has the potential to be worth between $3.9 and $11.1 trillion by 2025[@mck].

# Team
Datapace team is composed of experienced industry professionals, that together bring over 100 years of expertise in designing, building, deploying and maintaining large-scale industry networks.

**Drasko DRASKOVIC** - Drasko is an IoT expert with over 15 years of professional experience. He hacked embedded Linux SW and HW device drivers, designing complex wireless systems in telecom industry. Drasko earned his reputation in open-source community as an author of numerous projects - like WeIO[^weio] or Mainflux. He is one of the main contributors of the Linux Foundation's EdgeX Foundry[^edgex] project. Drasko is also author of the book "Scalable Architectures for the Internet of Things" published by O'Reilly and a vivid conference speaker. He holds a MSc in Electrical Engineering from Belgrade University.

**George SALEH** - George has strong expertise in IoT, ML and AI, with over 18 years of professional experience. Over the years, he gained significant knowledge in wireless management systems, designing complex OSS application and services in telecom industry. George earned his reputation as an entrepreneur of numerous ventures - like FreeFlex or IoT Gateway. He is also part of the incubation, innovation and maker place in Nokia, Paris - “Le Garage”. He holds a MSc in Computer Networking and Telecommunication from Paul Sabatier University, Toulouse.

**Rastko BLAGOJEVIC** - Rastko brings over 15 years of strong technical and management experience. Through his work in Siemens Networks and Nokia he has been exposed to complex projects in telecommunication industry and acquired valuable contacts working with biggest telecom operators and vendors in Europe and USA. Rastko holds MSc in Telecommunications and also MSc in Economy, both from Belgrade University.

**Janko ISIDOROVIC** - Janko is the co-founder of Mainflux IoT open source project. He is also chair of the Application Work Group of Linux Foundation EdgeX Foundry project. Janko has a background in Project Management, IT and Software integrations. He holds MSc in Telecommunications from Belgrade University and brings more than 15 years of extensive industry experience to Datapace project.

**Nikola MARCETIC** - Nikola is software developer with over 10 years of expertise in Web development and connected things over IP, delivering a complex and scalable Web apps. He enjoys in cutting edge technologies and projects with Linux OS & related open source tools with technical barriers and complex challenges. Nikola holds MSc in Economy from Novi Sad University.

**Darko DRASKOVIC** - Darko holds a PhD thesis on AI from UNIL, Switzerland. His special fields of interest include web development, graphics programming and data science.
Darko is also philosopher with a special interest in contemporary science and cutting edge IT technology. Darko has been hacking mathematically complex software solutions for over 10 years.

**Manuel IMPERIALE** - Manuel holds M.Sc. form University Pierre and Marie Curie in Paris. He has specialized in industrial informatics and both software and hardware technologies. Manuel was working in The Institute for Intelligent Systems and Robotics (ISIR), and various hi-tech companies, like 3D Sound Labs and Devialet.

**Sasa KLOPANOVIC** - Recently engaged in the biggest real estate project in South-East Europe master-planned by world renowned architects, Sasa gained significant working experience within the international environment where communication with international companies, eminent consultancies, and government institutions was a major part of his daily tasks. Additionally, he is the founder of first Serbian Branding Association and part of the team responsible for the development of National Brand of Serbia visual identity and slogan. He holds MSc degree in Philosophy from Belgrade University.

[^weio]: WeIO is an innovative open source hardware and software platform for rapid prototyping and creation of wirelessly connected interactive objects using only popular web languages such as HTML5 or Python. More information can be obtained at project website: <http://we-io.net>.

[^edgex]: EdgeX Foundry is a vendor-neutral open source project building a common open framework for IoT edge computing. At the heart of the project is an interoperability framework hosted within a full hardware- and OS-agnostic reference software platform to enable an ecosystem of plug-and-play components that unifies the marketplace and accelerates the deployment of IoT solutions. More information can be obtained at project's web address: <https://www.edgexfoundry.org/>.

# Contact
Website: <https://www.datapace.com>

E-mail: <info@datapace.com>

## Social Networks
Twitter: [\@DatapaceMarket](https://twitter.com/DatapaceMarket)

LinkedIn: <https://www.linkedin.com/company/datapace/>

Facebook: <https://www.facebook.com/datapace>

# Acknowledgments {-}
This work is the cumulative effort of multiple individuals within the Datapace team, and would not
have been possible without the help, comments, and review of the collaborators and advisors of Datapace. Drasko Draskovic ad George Saleh wrote the original Datapace whitepaper in 2017, laying the groundwork for this work. Nikola Marcetic developed ideas related to the Datapace decentralized protocol and implemented original temporary proxying scheme. Janko Isidorovic contributed ideas related to TAS token economy and integration with Mainflux IoT platform. Manuel imperiale, Darko Draskovic and Sasa Klopanovic improved the protocol, designed TAS token wallet and added geofencing capabilities. We also thank all of our collaborators and advisors for useful conversations; in particular Ilia Zelenikin and Dejan Novakovic.

# References
